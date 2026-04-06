package di

import "net/http"

type Server struct {
	Router http.Handler
	Port   string
}

func NewServer(router http.Handler, port string) *Server {
	return &Server{
		Router: router,
		Port:   port,
	}
}
