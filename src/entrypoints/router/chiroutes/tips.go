package chiroutes

import (
	"net/http"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/controllers"
	"github.com/go-chi/chi"
)

func Tips(ctrl *controllers.TipsController) http.Handler {
	// TODO: if trying to access a route that doesn't exist, 405 is returned by default but should it be 404?
	r := chi.NewRouter()
	r.Get("/", ctrl.GetAll)
	r.Post("/", ctrl.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ctrl.Get)
		r.Delete("/", ctrl.Delete)
	})
	return r
}
