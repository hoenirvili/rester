// Package token offers a basic jwt token validation method
package token

import (
	"crypto/rsa"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/hoenirvili/rester/permission"
)

type JWT struct {
	extractor request.Extractor
	keyFunc   jwt.Keyfunc
	options   []request.ParseFromRequestOption
	claims    *Claims
}

type Claims struct {
	mapClaims jwt.MapClaims
}

func (c *Claims) VerifyPermissions() error {
	p, ok := c.mapClaims["permissions"]
	if !ok {
		return errors.New("No permission found in the jwt token")
	}
	v := permission.Permissions(p.(float64))
	if !v.Valid() {
		return errors.New("Invalid permissions value, value not supported")
	}
	return nil
}

func (c *Claims) Valid() error {
	_, ok := c.mapClaims["exp"]
	if !ok {
		return errors.New("Jwt token exp field not present")
	}

	now := time.Now().Unix()
	if !c.mapClaims.VerifyExpiresAt(now, false) {
		return &jwt.ValidationError{
			Inner:  errors.New("Token is expired"),
			Errors: jwt.ValidationErrorExpired,
		}
	}
	if err := c.VerifyPermissions(); err != nil {
		return &jwt.ValidationError{Inner: err}
	}
	return nil
}

func NewClaims() *Claims {
	return &Claims{make(jwt.MapClaims)}
}
func NewJWT(key *rsa.PublicKey) *JWT {
	claims := NewClaims()
	opt := request.WithClaims(claims.mapClaims)
	fn := func(t *jwt.Token) (interface{}, error) { return key, nil }
	ext := request.AuthorizationHeaderExtractor
	return &JWT{ext, fn, []request.ParseFromRequestOption{opt}, claims}
}

func (j *JWT) Verify(r *http.Request) (map[string]interface{}, error) {
	t, err := request.ParseFromRequest(r,
		j.extractor, j.keyFunc, j.options...)
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("Jwt token is not valid")
	}
	return j.claims.mapClaims, nil
}
