package security

import (
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/config"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type SecurityClient interface {
	HashPassword(p string) (string, error)
	VerifyPassword(hp, p string) error
	GenerateJWTToken(uuid string) (string, error)
}

type client struct {
	securityConfig config.SecurityConfig
}

func NewSecurityClient(sc config.SecurityConfig) SecurityClient {
	return &client{sc}
}

func (c *client) HashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Stack(err)
	}

	return string(bytes), nil
}

func (c *client) VerifyPassword(hp, p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hp), []byte(p))
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}

func (c *client) GenerateJWTToken(uuid string) (string, error) {
	signature, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: c.securityConfig.HMAC512Key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", errors.Wrap(err, "signature")
	}

	cl := jwt.Claims{
		Subject: "login",
		Issuer:  "aymen",
	}

	raw, err := jwt.Signed(signature).Claims(cl).CompactSerialize()
	if err != nil {
		return "", errors.Wrap(err, "final compact")
	}
	raw2, err := jwt.Signed(signature).Claims(cl).FullSerialize()
	if err != nil {
		return "", errors.Wrap(err, "final full")
	}

	zap.S().Infof("raw = %v and raw2 full = %v", raw, raw2)

	return raw, nil
}
