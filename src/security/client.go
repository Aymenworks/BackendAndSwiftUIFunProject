package security

import (
	"net/http"
	"strings"
	"time"

	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/config"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/errors"
	"github.com/aymenworks/ProjectCookingTips-GoFromScratch/src/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type SecurityClient interface {
	HashPassword(p string) (string, error)
	VerifyPassword(hp, p string) error
	GenerateJWTToken(expiry time.Time, userUUID string) (*JWTToken, error)
	VerifyJWTTokenFromRequest(r *http.Request) (*PrivateClaimsJWT, error)
	VerifyJWTTokenFromString(t string) (*PrivateClaimsJWT, error)
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

type JWTToken struct {
	Token    string
	UUID     string
	UserUUID string
	Expiry   time.Time
}

func (c *client) GenerateJWTToken(expiry time.Time, userUUID string) (*JWTToken, error) {
	enc, err := jose.NewEncrypter(
		jose.A128GCM, // content encryption algo
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       []byte(c.securityConfig.JWTEncryptionKey),
		},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"))
	if err != nil {
		return nil, errors.Stack(err)
	}
	signature, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: []byte(c.securityConfig.JWTSignatureHMAC512Key)}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return nil, errors.Stack(err)
	}

	cl := jwt.Claims{
		Subject: "login",
		Issuer:  "aymen",
		Expiry:  jwt.NewNumericDate(expiry),
	}

	// TODO: set the struc somewhere to avoid having to guess the parameter when decoding
	pcl := &PrivateClaimsJWT{
		UUID:     uuid.NewString(),
		UserUUID: userUUID,
	}

	raw, err := jwt.SignedAndEncrypted(signature, enc).Claims(cl).Claims(pcl).CompactSerialize()
	if err != nil {
		return nil, errors.Stack(err)
	}

	token := &JWTToken{
		Token:    raw,
		UUID:     pcl.UUID,
		UserUUID: pcl.UserUUID,
		Expiry:   expiry,
	}

	return token, nil
}

type PrivateClaimsJWT struct {
	UUID     string `json:"uuid"`
	UserUUID string `json:"user_uuid"`
}

func (c *client) verifyJSONWebToken(ts string) (*PrivateClaimsJWT, error) {
	nestedTkn, err := jwt.ParseSignedAndEncrypted(ts)
	if err != nil {
		return nil, errors.Wrap(errors.TokenInvalid, err.Error())
	}

	t, err := nestedTkn.Decrypt([]byte(c.securityConfig.JWTEncryptionKey))
	if err != nil {
		return nil, errors.Wrap(errors.TokenInvalid, err.Error())
	}

	if !c.containsAlgorithm(t.Headers, jose.HS512) {
		return nil, errors.Wrap(errors.TokenInvalid, "algorithm invalid")
	}

	cl := jwt.Claims{}
	pcl := new(PrivateClaimsJWT)

	if err := t.Claims([]byte(c.securityConfig.JWTSignatureHMAC512Key), &cl, &pcl); err != nil {
		return nil, errors.Wrap(errors.TokenInvalid, err.Error())
	}

	// TODO: test expiry that is correctly validated
	if err = cl.ValidateWithLeeway(jwt.Expected{
		Subject: "login",
		Issuer:  "aymen",
		Time:    time.Now(),
	}, 0); err != nil {
		return nil, errors.Wrap(errors.TokenInvalid, err.Error())
	}

	if utils.IsEmpty(pcl.UUID) || utils.IsEmpty(pcl.UserUUID) {
		return nil, errors.TokenInvalid
	}

	return pcl, nil
}

func (c *client) containsAlgorithm(hds []jose.Header, alg jose.SignatureAlgorithm) bool {
	for _, o := range hds {
		if jose.SignatureAlgorithm(o.Algorithm) == alg {
			return true
		}
	}
	return false
}

func (c *client) extractToken(r *http.Request) string {
	tk := r.Header.Get("Authorization")
	if utils.IsEmpty(tk) {
		return ""
	}
	splt := strings.Split(tk, " ")
	if len(splt) != 2 {
		return ""
	}
	return splt[1]
}

func (c *client) VerifyJWTTokenFromRequest(r *http.Request) (*PrivateClaimsJWT, error) {
	et := c.extractToken(r)
	if utils.IsEmpty(et) {
		return nil, errors.TokenNotSet
	}
	pcl, err := c.verifyJSONWebToken(et)
	if err != nil {
		return nil, errors.Stack(err)
	}
	return pcl, nil
}

func (c *client) VerifyJWTTokenFromString(t string) (*PrivateClaimsJWT, error) {
	pcl, err := c.verifyJSONWebToken(t)
	if err != nil {
		return nil, errors.Stack(err)
	}
	return pcl, nil
}
