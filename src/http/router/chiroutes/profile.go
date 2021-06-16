package chiroutes

import (
	"net/http"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/controllers"
	"github.com/go-chi/chi"
)

func Profile(ctrl *controllers.ProfileController) http.Handler {
	// TODO: if trying to access a route that doesn't exist, 405 is returned by default but should it be 404?
	r := chi.NewRouter()
	r.Get("/", ctrl.Get)
	return r
}
