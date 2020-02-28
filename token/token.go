// Package token offers a basic jwt token validation method
package token

import (
	"crypto/rsa"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type JWT struct {
	extractor request.Extractor
	keyFunc   jwt.Keyfunc
	options   []request.ParseFromRequestOption
	claims    *Claims
}

func NewJWT(key *rsa.PublicKey) *JWT {
	claims := NewClaims()
	opt := request.WithClaims(claims)
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
		return nil, errors.New("jwt token is not valid")
	}
	return j.claims.mapClaims, nil
}
