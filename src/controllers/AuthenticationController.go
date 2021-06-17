package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/api"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/services/user"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	requests "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/http/requests/tips"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/http/responses"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/infra/caches"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/security"
	"go.uber.org/zap"
)

type AuthenticationController struct {
	Controller
	userSvc        user.Service
	securityClient security.SecurityClient
	cacheClt       caches.Cache
}

func NewAuthenticationController(userSvc user.Service, securityClient security.SecurityClient, cacheClt caches.Cache) *AuthenticationController {
	return &AuthenticationController{
		userSvc:        userSvc,
		securityClient: securityClient,
		cacheClt:       cacheClt,
	}
}

func (c *AuthenticationController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	acRtUUID, err := c.cacheClt.Get(r.Context(), fmt.Sprintf("access_token:%v", api.AccessTokenUUID(ctx)))
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}

	// Delete refresh token
	acRtUUIDString, ok := acRtUUID.(string)
	if ok {
		if err := c.cacheClt.Delete(ctx, fmt.Sprintf("refresh_token:%v", acRtUUIDString)); err != nil {
			c.ErrorResponse(w, errors.Stack(err))
			return
		}
	} else {
		zap.S().Errorf("Could not delete the refresh token: cast to string error")
	}

	// delete access token
	if err := c.cacheClt.Delete(ctx, fmt.Sprintf("access_token:%v", api.AccessTokenUUID(ctx))); err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	c.NoContentResponse(w)
}

func (c *AuthenticationController) Login(w http.ResponseWriter, r *http.Request) {
	var req requests.LoginRequest
	if err := c.ParseBody(r, &req); err != nil {
		c.ErrorResponse(w, err)
		return
	}

	u, err := c.userSvc.MustGetByUsername(r.Context(), req.Username)
	if err != nil {
		// TODO: replace with appropriate error in case db return something important like a critical error
		c.ErrorResponse(w, errors.Stack(errors.IncorrectCredentials))
		return
	}

	err = c.securityClient.VerifyPassword(u.HashedPassword, req.Password)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(errors.IncorrectCredentials))
		return
	}

	// TODO: Try to encrypt UUID for private I guess
	now := time.Now()
	at, err := c.securityClient.GenerateJWTToken(now.Add(time.Second*20), u.UUID)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}
	rt, err := c.securityClient.GenerateJWTToken(now.Add(time.Hour*24*7), u.UUID)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	if err = c.cacheClt.Set(r.Context(), fmt.Sprintf("access_token:%v", at.UUID), rt.UUID, at.Expiry.Sub(now)); err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}
	if err = c.cacheClt.Set(r.Context(), fmt.Sprintf("refresh_token:%v", rt.UUID), at.UUID, rt.Expiry.Sub(now)); err != nil {
		_ = c.cacheClt.Delete(r.Context(), fmt.Sprintf("access_token:%v", at.UUID))
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	// TODO: CHeck KMS for storing keys and have APIs -> available on LocalStack

	res := &responses.SessionToken{
		AccessToken:  at.Token,
		RefreshToken: rt.Token,
	}

	c.JsonResponse(w, res)
}

func (c *AuthenticationController) Signup(w http.ResponseWriter, r *http.Request) {
	var req requests.SignUpRequest
	if err := c.ParseBody(r, &req); err != nil {
		c.ErrorResponse(w, err)
		return
	}

	u, err := c.userSvc.GetByUsername(r.Context(), req.Username)
	if err != nil {
		// TODO: handle appriopriately the error
		c.ErrorResponse(w, err)
		return
	}

	if u != nil {
		c.ErrorResponse(w, errors.UsernameAlreadyAssigned)
		return
	}

	hp, err := c.securityClient.HashPassword(req.Password)
	if err != nil {
		// TODO: handle appriopriately the error
		c.ErrorResponse(w, err)
		return
	}
	_, err = c.userSvc.Create(r.Context(), req.Username, hp)
	if err != nil {
		// TODO: handle appriopriately the error
		c.ErrorResponse(w, err)
		return
	}

	// TODO: check which status is usually returned for signing up or maybe the user can automatically be logged in
	c.NoContentResponse(w)
}

func (c *AuthenticationController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req requests.RefreshToken
	if err := c.ParseBody(r, &req); err != nil {
		c.ErrorResponse(w, err)
		return
	}
	pcl, err := c.securityClient.VerifyJWTTokenFromString(req.Token)
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}
	rtAcUUID, err := c.cacheClt.Get(r.Context(), fmt.Sprintf("refresh_token:%v", pcl.UUID))
	if err != nil {
		c.ErrorResponse(w, err)
		return
	}
	if rtAcUUID == nil {
		c.ErrorResponse(w, errors.TokenNotFound)
		return
	}
	// Create new tokens
	// TODO: Try to encrypt UUID for private I guess
	zap.S().Debugf("pcl.UserUUID = %v", pcl.UserUUID)
	now := time.Now()
	at, err := c.securityClient.GenerateJWTToken(now.Add(time.Minute*15), pcl.UserUUID)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}
	rt, err := c.securityClient.GenerateJWTToken(now.Add(time.Hour*24*7), pcl.UserUUID)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}
	// Delete old access token
	ctx := r.Context()
	rtAcUUIDString, ok := rtAcUUID.(string)
	if ok {
		zap.S().Debugf("rtAcUUIDString = %v", rtAcUUIDString)
		if err := c.cacheClt.Delete(ctx, fmt.Sprintf("access_token:%v", rtAcUUIDString)); err != nil {
			c.ErrorResponse(w, errors.Stack(err))
			return
		}
	}
	// delete old refresh token
	if err := c.cacheClt.Delete(ctx, fmt.Sprintf("refresh_token:%v", pcl.UUID)); err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	// Save the new tokens
	if err = c.cacheClt.Set(r.Context(), fmt.Sprintf("access_token:%v", at.UUID), rt.UUID, at.Expiry.Sub(now)); err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}
	if err = c.cacheClt.Set(r.Context(), fmt.Sprintf("refresh_token:%v", rt.UUID), at.UUID, rt.Expiry.Sub(now)); err != nil {
		_ = c.cacheClt.Delete(r.Context(), fmt.Sprintf("access_token:%v", at.UUID))
		c.ErrorResponse(w, errors.Stack(err))
		return
	}

	res := &responses.SessionToken{
		AccessToken:  at.Token,
		RefreshToken: rt.Token,
	}

	c.JsonResponse(w, res)
}
