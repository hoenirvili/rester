// Package handler used for defining custom handler types used
// by the internal HTTP router
package handler

import (
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
)

type Handler func(request.Request) resource.Response
