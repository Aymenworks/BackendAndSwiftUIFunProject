package middlewares

import (
	"context"
	"net/http"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/api"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/security"
)

func AuthenticatedOnly(secClient security.SecurityClient, verifyAccessToken func(ctx context.Context, uuid string) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pcl, err := secClient.VerifyJWTToken(r)
			if err != nil {
				api.ErrorResponse(w, errors.Wrap(errors.TokenInvalid, err.Error()))
				return
			}
			ctx := r.Context()
			ctx = api.WithUserUUID(ctx, pcl.UserUUID)
			ctx = api.WithAccessTokenUUID(ctx, pcl.UUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
