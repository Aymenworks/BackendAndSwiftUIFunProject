package chiroutes

import (
	"net/http"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/controllers"
	"github.com/go-chi/chi"
)

func Auth(ctrl *controllers.AuthenticationController) http.Handler {
	// TODO: if trying to access a route that doesn't exist, 405 is returned by default but should it be 404?
	r := chi.NewRouter()
	r.Post("/signup", ctrl.Signup)
	r.Post("/login", ctrl.Login)

	// TODO: requires token
	r.Get("/refresh-token", ctrl.RefreshToken)
	// TODO: requires token
	r.Post("/logout", ctrl.RefreshToken)

	return r
}
