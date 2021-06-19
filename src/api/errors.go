package api

import (
	"encoding/json"
	"net/http"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"go.uber.org/zap"
)

func ErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	zap.S().Error(err)
	apperr := errors.AsAppError(err)
	w.WriteHeader(apperr.HTTPStatusCode)

	if apperr.CanDisplayMessage() {
		if err := json.NewEncoder(w).Encode(apperr.New()); err != nil {
			zap.S().Errorf("Error encoding error response %v", err)
			return
		}
	}
}
