package controllers

import (
	"net/http"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/domain/services/user"
	requests "github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/entrypoints/requests/tips"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/security"
)

type AuthenticationController struct {
	Controller
	userSvc        user.Service
	securityClient security.SecurityClient
}

func NewAuthenticationController(userSvc user.Service, securityClient security.SecurityClient) *AuthenticationController {
	return &AuthenticationController{
		userSvc:        userSvc,
		securityClient: securityClient,
	}
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

	t, err := c.securityClient.GenerateJWTToken(u.UUID)
	if err != nil {
		c.ErrorResponse(w, errors.Stack(err))
		return
	}
	// TODO:save the token in redis
	c.JsonResponse(w, t)
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

	// TODO:
	// - use the Login service after having parsed some credentials, ideally hashed as well on the client side.
	// - check with the hashed credentials in the db whether it match or return StatusUnauthorized with least detail possible on the error message
	// - if OK, create the token, save the token in redis, return it to the user
	c.NoContentResponse(w)
}
