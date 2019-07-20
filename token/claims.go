package token

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hoenirvili/rester/permission"
)

type Claims struct {
	mapClaims jwt.MapClaims
}

func NewClaims() *Claims {
	return &Claims{make(jwt.MapClaims)}
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

func (c *Claims) VerifyExp() error {
	_, ok := c.mapClaims["exp"]
	if !ok {
		return errors.New("Jwt token exp field not present")
	}
	return nil
}

func (c *Claims) Valid() error {
	now := time.Now().Unix()
	if !c.mapClaims.VerifyExpiresAt(now, true) {
		return &jwt.ValidationError{
			Inner:  errors.New("Token is expired"),
			Errors: jwt.ValidationErrorExpired,
		}
	}
	return nil
}

func (c *Claims) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &c.mapClaims); err != nil {
		return err
	}
	if err := c.VerifyPermissions(); err != nil {
		return &json.UnsupportedValueError{
			Value: reflect.ValueOf(c.mapClaims["permissions"]),
			Str:   err.Error(),
		}
	}
	if err := c.VerifyExp(); err != nil {
		return &json.UnsupportedValueError{
			Value: reflect.ValueOf(c.mapClaims["exp"]),
			Str:   err.Error(),
		}
	}
	return nil
}

var (
	_ json.Unmarshaler = (*Claims)(nil)
	_ jwt.Claims       = (*Claims)(nil)
)
