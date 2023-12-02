package webapi

import (
	"net/http"
)

type IServer interface {
	RunServer(port string, handler http.Handler) error
}
type Server struct {
	httpServer *http.Server
}

func (s *Server) RunServer(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
	}
	return s.httpServer.ListenAndServe()
}

func InitServer() IServer {
	return &Server{}
}
