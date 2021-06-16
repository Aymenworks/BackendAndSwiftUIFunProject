package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/api"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	requests "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/http/requests/tips"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
	"github.com/go-chi/chi"
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
	api.ErrorResponse(w, err)
}

func (c *Controller) ParseContentTypeHeader(r *http.Request) (string, error) {
	ct := r.Header.Get("Content-type")
	t, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("error parsing content media type %v", ct))
	}

	return t, nil
}

func (c *Controller) ParseBody(r *http.Request, req requests.AppRequest) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		return errors.Wrap(err, "error parse body")
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
