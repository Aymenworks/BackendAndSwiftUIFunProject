package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/config"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/http/middlewares"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/http/router"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func main() {
	// Load conf
	config := config.New()

	// Setup router and middlewares
	router := router.NewChiRouter()

	// Log each request
	router.UseMiddleware(middleware.DefaultLogger)

	// TODO: Panic recover not working
	router.UseMiddleware(middlewares.PanicRecover)

	// e.g `/hello` and `/hello/` will be the same
	router.UseMiddleware(middleware.StripSlashes)

	// TODO: Confirm that it just stop readings after reaching that max
	router.UseMiddleware(middlewares.RequestSizeLimit)

	// TODO: Check middleware for response compression after having some rest APIs working

	// cancel a request if processing takes longer than 60 seconds,
	// server will respond with a http.StatusGatewayTimeout (504).
	// TODO: Not timing out unless using ctx-Done()
	router.UseMiddleware(middleware.Timeout(config.Middleware.Timeout))

	// TODO: Instanciate development or production based on environment
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // TODO: Understand why it's required to do
	sugar := logger.Sugar()

	router.Get("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}))

	// router.Get("/slow-with-ctx-done", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	sugar.Info("Start slow request\n")
	// 	ctx := r.Context()
	// 	processTime := time.Duration(7) * time.Second
	// 	select {
	// 	case <-ctx.Done():
	// 		sugar.Infof("Context Done() / err %v \n", ctx.Err())
	// 		return
	// 	case <-time.After(processTime):
	// 		sugar.Info("Request process after 7 seconds\n")
	// 	}

	// 	sugar.Info("Finish slow request\n")
	// 	w.Write([]byte("done"))
	// }))

	// Run the server in a goroutine to avoid blocking the current thread
	// so that we could allow listen for the shutdown signal right after
	server := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           router,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
		ReadTimeout:       config.Server.ReadTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			sugar.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Graceful shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		os.Interrupt,
	)

	sugar.Infof("received shutdown signal %v\n", <-shutdownCh)
	// Use a context with a timeout in case the Shutdown takes too much time and block the process
	// Canceling this context releases resources associated with it
	ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout)
	defer cancel()

	sugar.Infof("Shutting down\n")
	if err := server.Shutdown(ctx); err != nil {
		sugar.Errorf("Error shutting down %v\n", err)
	}
}
