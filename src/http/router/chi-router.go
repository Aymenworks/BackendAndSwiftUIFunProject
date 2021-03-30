package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

type router struct{}

func NewChiRouter() Router {
	return &router{}
}

var (
	chiRouter = chi.NewRouter()
)

func (*router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chiRouter.ServeHTTP(w, r)
}

func (*router) Get(uri string, f http.HandlerFunc) {
	chiRouter.Get(uri, f)
}

func (*router) Post(uri string, f http.HandlerFunc) {
	chiRouter.Post(uri, f)
}

func (*router) Put(uri string, f http.HandlerFunc) {
	chiRouter.Put(uri, f)
}

func (*router) Delete(uri string, f http.HandlerFunc) {
	chiRouter.Delete(uri, f)
}

func (*router) UseMiddleware(f func(http.Handler) http.Handler) {
	chiRouter.Use(f)
}

func (*router) Mount(uri string, h http.Handler) {
	chiRouter.Mount(uri, h)
}
