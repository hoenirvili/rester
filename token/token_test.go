package token_test

import (
	"testing"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/hoenirvili/rester/token"
	"github.com/stretchr/testify/require"
)

func TestNewJWT(t *testing.T) {
	require := require.New(t)
	token := token.NewJWT(nil)
	require.NotEmpty(token)
}

func jwt(r *require.Assertions) *token.JWT {
	public, err := gojwt.ParseRSAPublicKeyFromPEM(pub)
	r.NoError(err)
	token := token.NewJWT(public)
	r.NotEmpty(token)
	return token
}

func TestVerifyToken(t *testing.T) {
	require := require.New(t)
	_ = require
}

// type jwtSuite struct {
// 	suite.Suite
// 	token *token.JWT
// }
//
// func (s *jwtSuite) SetupTest() {
// 	require := s.Require()
// 	jwt := token.NewJWT(nil)
// 	require.NotEmpty(jwt)
//
// }
//
// func (s *jwtSuite) TearDownTest() {
//
// }
//
// func TestJWTToken(t *testing.T) {
// 	suite.Run(t, new(jwtSuite))
// }
