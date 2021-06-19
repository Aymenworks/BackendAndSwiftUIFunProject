package requests

import (
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *SignUpRequest) Validate() error {
	if utils.IsEmpty(r.Username) {
		return errors.UserNameEmpty
	}
	if utils.IsEmpty(r.Password) {
		return errors.PasswordEmpty
	}
	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if utils.IsEmpty(r.Username) {
		return errors.UserNameEmpty
	}
	if utils.IsEmpty(r.Password) {
		return errors.PasswordEmpty
	}
	return nil
}

type RefreshToken struct {
	Token string `json:"token"`
}

func (r *RefreshToken) Validate() error {
	if utils.IsEmpty(r.Token) {
		return errors.RefreshTokenEmpty
	}
	return nil
}
