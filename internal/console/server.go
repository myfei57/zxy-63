package console

import (
	"net/http"

	"waterplant/internal/store"
)

type Server struct {
	rt     *Runtime
	router http.Handler
}

func NewServer(s *store.Store) *Server {
	server := &Server{rt: NewRuntime(s)}
	server.router = server.routes()
	return server
}

func (s *Server) Handler() http.Handler {
	return s.router
}
