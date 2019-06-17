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
	fn := func(*jwt.Token) (interface{}, error) { return key, nil }
	ext := request.AuthorizationHeaderExtractor
	return &JWT{ext, fn, []request.ParseFromRequestOption{opt}, claims}
}

func (j JWT) Verify(r *http.Request) error {
	jwt, err := request.ParseFromRequest(r,
		j.extractor, j.keyFunc, j.options...)
	if err != nil {
		return err
	}

	if !jwt.Valid {
		return errors.New("JWT token is not valid")
	}

	return nil
}

func (j JWT) Extract() (permission.Permissions, error) {
	p, ok := j.claims["permissions"]
	if !ok {
		return 0, errors.New("No permission found in the jwt token")
	}
	v, ok := p.(permission.Permissions)
	if !ok {
		return v, errors.New("Invalid permissions value format")
	}
	return v, nil
}
