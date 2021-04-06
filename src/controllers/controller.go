package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	requests "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/entrypoints/requests/tips"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
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

	// TODO add context there looks at fields

	apperr := errors.AsAppError(err)
	zap.S().Error(err)
	w.WriteHeader(apperr.HTTPStatusCode)

	if err := json.NewEncoder(w).Encode(apperr.New()); err != nil {
		zap.S().Errorf("Error encoding error response %v", err)
		return
	}
}

func (c *Controller) ParseBody(r *http.Request, req requests.AppRequest) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		return errors.Wrap(errors.InvalidParameter, "error parse body")
	}

	err := req.Validate()
	if err != nil {
		return errors.Wrap(err, "error validate request")
	}

	return nil
}

func (c *Controller) PathParameterUint(r *http.Request, key string) (uint, error) {
	paramStr := chi.URLParam(r, key)
	if utils.IsEmpty(paramStr) {
		return 0, errors.Wrap(errors.PathKeyInvalid, fmt.Sprintf("cannot find the path key %v", key))
	}

	param, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		return 0, errors.Wrap(errors.PathKeyInvalid, fmt.Sprintf("cannot parse the id = %v into int", param))
	}

	return uint(param), nil
}
