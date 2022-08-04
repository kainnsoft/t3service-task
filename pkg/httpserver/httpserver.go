package httpserver

import (
	"context"
	"net/http"
	"team3-task/config"
	"team3-task/pkg/logging"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":3000"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New -.
func New(handler http.Handler, cfg config.HTTP, log *logging.ZeroLogger) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
		//ErrorLog:     log.ErrorLog, // вопрос ментору: как лучше передавать сюда эту зависимость? или лучше здесь использовать стандартный логгер?
	}

	if time.Duration(cfg.HttpReadTimeout) != _defaultReadTimeout {
		httpServer.ReadTimeout = time.Duration(cfg.HttpReadTimeout * int(time.Second))
	}

	if time.Duration(cfg.HttpWriteTimeout) != _defaultWriteTimeout {
		httpServer.WriteTimeout = time.Duration(cfg.HttpWriteTimeout * int(time.Second))
	}

	if cfg.HttpAddr != _defaultAddr {
		httpServer.Addr = cfg.HttpAddr
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) GetAddr() string {
	return s.server.Addr
}
