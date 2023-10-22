package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

type Server struct {
	router *mux.Router
	logger *slog.Logger
	host   string
	port   string
}

type Logger interface {
	// TODO
}

type Application interface { // TODO
}

func NewServer(logger *slog.Logger, host, port string, app Application) *Server { //nolint:revive
	return &Server{
		logger: logger,
		router: mux.NewRouter(),
		host:   host,
		port:   port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.configureRouter()
	serverAdd := fmt.Sprintf("%s:%s", s.host, s.port)
	if err := http.ListenAndServe(serverAdd, s.router); err != nil { //nolint:gosec
		s.logger.Info("failed to start server")
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error { //nolint:revive
	// TODO
	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/hello", loggingMiddleware(s.handleHello(), s.logger))
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
		w.WriteHeader(http.StatusOK)
	}
}
