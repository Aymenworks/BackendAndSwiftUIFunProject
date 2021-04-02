package controllers

import (
	"net/http"
	"strconv"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/services/tips"
	requests "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/entrypoints/requests/tips"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type TipsController struct {
	Controller
	service tips.Service
}

func NewTipsController(service tips.Service) *TipsController {
	return &TipsController{
		service: service,
	}
}

func (c *TipsController) GetAll(w http.ResponseWriter, r *http.Request) {
	tips, err := c.service.GetAll(r.Context())
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}

	c.JsonResponse(w, tips)
}

func (c *TipsController) Get(w http.ResponseWriter, r *http.Request) {
	tipID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}

	zap.S().Debugf("tipID = %v", tipID)
	tip, err := c.service.GetByID(r.Context(), uint(tipID))
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}

	c.JsonResponse(w, tip)
}

func (c *TipsController) Create(w http.ResponseWriter, r *http.Request) {
	var request requests.CreateTipRequest
	if err := c.ParseBody(r, &request); err != nil {
		c.ErrorResponse(w, err)
		return
	}
	zap.S().Debugf("request.name = %v", request.Name)

	c.NoContentResponse(w)
}

func (c *TipsController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := c.PathParameterUint(r, "id")
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}
	err = c.service.DeleteByID(id)
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}

	c.NoContentResponse(w)
}
