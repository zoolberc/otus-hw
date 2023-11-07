package internalhttp

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

type Server struct {
	router *mux.Router
	logger *slog.Logger
	host   string
	port   string
	srv    *http.Server
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
	serverAdd := net.JoinHostPort(s.host, s.port)
	s.srv = &http.Server{ //nolint:gosec
		Handler: s.router,
		Addr:    serverAdd,
	}
	if err := s.srv.ListenAndServe(); err != nil {
		s.logger.Info("failed to start server")
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
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
