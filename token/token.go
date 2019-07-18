// Package token offers a basic jwt token validation method
package token

import (
	"crypto/rsa"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/hoenirvili/rester/permission"
)

type JWT struct {
	extractor request.Extractor
	keyFunc   jwt.Keyfunc
	options   []request.ParseFromRequestOption
	claims    jwt.MapClaims
}

func NewJWT(key *rsa.PublicKey) *JWT {
	claims := jwt.MapClaims{}
	opt := request.WithClaims(claims)
	fn := func(t *jwt.Token) (interface{}, error) {
		_, ok := claims["exp"]
		if !ok {
			return nil, errors.New("Jwt token exp field not present")
		}
		return key, nil
	}
	ext := request.AuthorizationHeaderExtractor
	return &JWT{ext, fn, []request.ParseFromRequestOption{opt}, claims}
}

func (j JWT) Verify(r *http.Request) error {
	t, err := request.ParseFromRequest(r,
		j.extractor, j.keyFunc, j.options...)
	if err != nil {
		return err
	}
	if !t.Valid {
		return errors.New("Jwt token is not valid")
	}
	return nil
}

func (j JWT) Extract() (permission.Permissions, error) {
	p, ok := j.claims["permissions"]
	if !ok {
		return 0, errors.New("No permission found in the jwt token")
	}
	v := permission.Permissions(p.(float64))
	if !v.Valid() {
		return v, errors.New("Invalid permissions value, value not supported")
	}
	return v, nil
}
