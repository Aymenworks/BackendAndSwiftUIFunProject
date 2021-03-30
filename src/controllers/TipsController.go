package controllers

import (
	"net/http"
	"strconv"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/services/tips"
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
	}

	c.JsonResponse(w, tips)
}

func (c *TipsController) Get(w http.ResponseWriter, r *http.Request) {
	tipID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		c.ErrorResponse(w, err)
	}

	zap.S().Debugf("tipID = %v", tipID)
	tip, err := c.service.GetByID(r.Context(), uint(tipID))
	if err != nil {
		c.ErrorResponse(w, err)
	}

	c.JsonResponse(w, tip)
}

func (c *TipsController) Update(w http.ResponseWriter, r *http.Request) {
	tips, err := c.service.GetAll(r.Context())
	if err != nil {
		c.ErrorResponse(w, err)
	}

	c.JsonResponse(w, tips)
}

func (c *TipsController) Delete(w http.ResponseWriter, r *http.Request) {
	tips, err := c.service.GetAll(r.Context())
	if err != nil {
		c.ErrorResponse(w, err)
	}

	c.JsonResponse(w, tips)
}
