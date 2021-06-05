package utils

import (
	"crypto"
	"crypto/hmac"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/config"
	"golang.org/x/crypto/bcrypt"
)

type client struct {
	SecurityConfig config.SecurityConfig
}

func NewSecurityClient(sc config.SecurityConfig) *client {
	return &client{sc}
}

func (c *client) HashPassword(p string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (c *client) ComparePassword(p string, h []byte) error {
	if err := bcrypt.CompareHashAndPassword(h, []byte(p)); err != nil {
		return err
	}
	return nil
}

var tmpKey []byte

func (c *client) SignMessage(msg []byte) ([]byte, error) {
	h := hmac.New(crypto.SHA512.New, []byte(c.SecurityConfig.HMACKey))
	_, err := h.Write(msg)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (c *client) CheckSignature(msg, s []byte) (bool, error) {
	newS, err := c.SignMessage(msg)
	if err != nil {
		return false, err
	}

	same := hmac.Equal(s, newS)
	return same, nil
}
