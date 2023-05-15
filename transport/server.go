package transport

import (
	"net/http"

	"test/service"
	"test/transport/handler"
)

type Server struct {
	handler *handler.Manager
	mux     *http.ServeMux
}

func NewServerHTTP(service *service.Manager) *Server {
	return &Server{
		handler: handler.NewManagerHandler(service),
		mux:     new(http.ServeMux),
	}
}

func (s *Server) Run() error {
	s.route()
	return http.ListenAndServe(":8080", s.mux)
}
