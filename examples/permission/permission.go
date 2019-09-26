package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/hoenirvili/rester"
	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
	"github.com/hoenirvili/rester/response"
	"github.com/hoenirvili/rester/route"
	"github.com/hoenirvili/rester/token"
)

type jsonResponse struct {
	Message string `json:"message"`
}

type root struct{}

func (r *root) index(req request.Request) resource.Response {
	return response.Payload(&jsonResponse{"Hello World !"})
}

func (r *root) admin(req request.Request) resource.Response {
	return response.Payload(&jsonResponse{"Hello Admin !"})
}

func (r *root) Routes() route.Routes {
	return route.Routes{
		// This will allow anyone per default
		{
			URL:     "/",
			Method:  resource.Get,
			Handler: r.index,
		},
		// Only admins
		{

			URL:     "/admin",
			Method:  resource.Get,
			Allow:   permission.Admin,
			Handler: r.admin,
		},
	}
}

var pub = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6X5gMNrebNra73WI7jIX
jQSTiBEvI79hMlT+b2qrxDKtq0wY8CaNkr9k5sFyzE/ImGo4NHz3rArkKIQU5uSy
Y0xjJvN0uhqmNGv3F+yXf/ak9BIyiq6TM45j15KVbw+dLHSZJhN8ZZ96C36dOyB2
JPInjkU+pMzJr/Dp6KBL5dDqYqSQzmojb0KxyNyn1VGvurG5e3TFCH4377bg+H0Z
C48AKkhokFvhhI7MJpFhfoQaGqAgg3kOguWOpngDcaTwschjwZ2rW9/qguf6iSsG
77+5RyJrlkZbbpaQl83gEs2EMzykmhzfhDJWXsEA0+HH4ns8XDgqlIodIkTmf2/5
lQIDAQAB
-----END PUBLIC KEY-----
`)[1:]

func main() {
	key, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		panic(err)
	}
	token := token.NewJWT(key)
	rester := rester.New(
		rester.WithTokenValidator(token),
	)
	rester.Resource("/", new(root))
	rester.Build()
	if err := http.ListenAndServe(":8080", rester); err != nil {
		panic(err)
	}
}
