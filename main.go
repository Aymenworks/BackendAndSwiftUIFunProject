package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/config"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/controllers"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/services/tips"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/entrypoints/router"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/entrypoints/router/chiroutes"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/http/middlewares"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/infra/databases"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/infra/mysql"
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

	// TODO: Check Panic recover not working
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
	appLog := initAppLog()
	zap.ReplaceGlobals(appLog)
	defer appLog.Sync() // TODO: Understand why it's required to do
	logger := appLog.Sugar()

	// Config DB
	db := databases.NewMySQLDatabase(config.DockerPort, config.Database)

	// Run the server in a goroutine to avoid blocking the current thread
	// so that we could allow listen for the shutdown signal right after
	// TODO: Confirm for the timeout stuff after some routes are implemented
	server := &http.Server{
		Addr:              ":" + config.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
		ReadTimeout:       config.Server.ReadTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
	}

	// Setup repositories
	tipsRepository := mysql.NewMysqlTipsRepository(db)

	// Setup services
	tipsService := tips.NewTipsService(tipsRepository)

	// Setup controllers
	profileController := controllers.NewTipsController(tipsService)

	// Setup routes
	router.Mount("/tips", chiroutes.Tips(profileController))

	go func() {
		logger.Debugf("HTTP server ListenAndServe: %v", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("Error HTTP server ListenAndServe: %v", err)
		}
	}()

	// Graceful shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		os.Interrupt,
	)

	logger.Infof("received shutdown signal %v\n", <-shutdownCh)
	// Use a context with a timeout in case the Shutdown takes too much time and block the process
	// Canceling this context releases resources associated with it
	ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout)
	defer cancel()

	logger.Infof("Shutting down\n")
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Error shutting down %v\n", err)
	}
}

func initAppLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	// TODO: Check if any config is useful
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Could not init zap loggeer with err = %v", err)
	}
	return logger
}
