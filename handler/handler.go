package handler

import (
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
)

type Handler func(request.Request) resource.Response
