package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	requests "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/entrypoints/requests/tips"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type Controller struct {
}

func (c *Controller) JsonResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		c.ErrorResponse(w, err)
		return
	}
}

func (c *Controller) NoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) ErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	appErr := errors.AsAppError(err)
	if appErr == nil {
		appErr = errors.NewUnknownError()
	}

	// TODO: To replace, of course
	w.WriteHeader(appErr.HTTPStatusCode)

	if err := json.NewEncoder(w).Encode(appErr); err != nil {
		zap.S().Errorf("Error encoding error response %w", err)
		return
	}
}

func (c *Controller) ParseBody(r *http.Request, req requests.AppRequest) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.S().Debugf("Error decoding body %v", err)
		return xerrors.Errorf("Error decoding body %v", err)
	}

	zap.S().Debugf("Decoding was ok with req = %+v", req)

	return req.Validate()
}

func (c *Controller) PathParameterUint(r *http.Request, key string) (uint, error) {
	paramStr := chi.URLParam(r, key)
	if utils.IsEmpty(paramStr) {
		// TODO: return appropriates errors
		return 0, xerrors.Errorf("Cannot find the path key %v", key)
	}
	param, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		// TODO: return appropriates errors
		return 0, xerrors.Errorf("Cannot parse the id = %v into int", param)
	}

	return uint(param), nil
}
