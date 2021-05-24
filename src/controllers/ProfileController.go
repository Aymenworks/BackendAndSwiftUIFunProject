package controllers

import (
	"net/http"
	"time"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/infra/caches"
	"go.uber.org/zap"
)

type ProfileController struct {
	Controller
	cacheClt caches.Cache
}

func NewProfileController(cacheClt caches.Cache) *ProfileController {
	return &ProfileController{
		cacheClt: cacheClt,
	}
}
func (c *ProfileController) Get(w http.ResponseWriter, r *http.Request) {
	cache, err := c.cacheClt.Get(r.Context(), "profile")
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}

	if cache != nil {
		zap.S().Info("Reuse cache")
		c.JsonResponse(w, cache)
		return
	}

	zap.S().Info("Set cache for the first time")
	intialVal := "initial-value"
	if err = c.cacheClt.Set(r.Context(), "profile", intialVal, time.Second*5); err != nil {
		c.ErrorResponse(w, err)
		return
	}

	c.JsonResponse(w, intialVal)
}
