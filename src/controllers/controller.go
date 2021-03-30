package controllers

import (
	"encoding/json"
	"net/http"
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
	// TODO: To replace, of course
	w.WriteHeader(http.StatusInternalServerError)
}
