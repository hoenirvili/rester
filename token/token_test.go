package token_test

import (
	"testing"
	"time"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"

	"github.com/hoenirvili/rester/token"
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

func tokenString(r *require.Assertions) string {
	private, err := gojwt.ParseRSAPrivateKeyFromPEM(priv)
	r.NoError(err)
	token := gojwt.NewWithClaims(gojwt.SigningMethodRS256, gojwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	stoken, err := token.SignedString(private)
	r.NoError(err)
	return stoken
}

func tokenStringWithoutAnyClaims(r *require.Assertions) string {
	private, err := gojwt.ParseRSAPrivateKeyFromPEM(priv)
	r.NoError(err)
	token := gojwt.New(gojwt.SigningMethodRS256)
	stoken, err := token.SignedString(private)
	r.NoError(err)
	return stoken
}

func tokenStringExpired(r *require.Assertions) string {
	private, err := gojwt.ParseRSAPrivateKeyFromPEM(priv)
	r.NoError(err)
	token := gojwt.NewWithClaims(gojwt.SigningMethodRS256, gojwt.MapClaims{
		"exp": time.Now().Add(time.Hour * -1).Unix(),
	})
	stoken, err := token.SignedString(private)
	r.NoError(err)
	return stoken
}

// func TestVerifyToken(t *testing.T) {
// 	require := require.New(t)
// 	token := jwt(require)
// 	tokenString := tokenString(require)
// 	_, err := token.Verify(&http.Request{
// 		Header: http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", tokenString)}},
// 	})
// 	require.NoError(err)
// }

// func TestVerifyTokenWithBadInput(t *testing.T) {
// 	require := require.New(t)
// 	inputs := []string{
// 		fmt.Sprintf("Bearer %s", "badtoken"),
// 		"Bearer", "", "sfhjiuasdhgiusdfsdfa",
// 		fmt.Sprintf("Bearer %s", tokenStringExpired(require)),
// 		fmt.Sprintf("Bearer %s", tokenStringWithoutAnyClaims(require)),
// 	}
// 	for _, input := range inputs {
// 		token := jwt(require)
// 		_, err := token.Verify(&http.Request{
// 			Header: http.Header{"Authorization": []string{input}},
// 		})
// 		require.Error(err)
// 	}
// }
