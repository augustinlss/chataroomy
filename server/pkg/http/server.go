package http

import "net/http"

type Server struct {
	Addr           string
	Handler        http.Handler
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	MaxHeaderBytes int
	MaxConnsPerIP  int
}
