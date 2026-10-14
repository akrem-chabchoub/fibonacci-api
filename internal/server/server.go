// Package server runs the HTTP API and shuts it down cleanly.
package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 10 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultIdleTimeout       = 60 * time.Second
)

// Options configure the HTTP server.
type Options struct {
	Addr              string
	ShutdownTimeout   time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration
	Logger            *zap.Logger
}

// Server is the process HTTP listener.
type Server struct {
	http            *http.Server
	log             *zap.Logger
	shutdownTimeout time.Duration
	mu              sync.Mutex
	bound           string
	started         chan struct{}
}

// New builds a server that is not listening yet.
func New(opts Options) *Server {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	if opts.ShutdownTimeout <= 0 {
		opts.ShutdownTimeout = 10 * time.Second
	}

	if opts.ReadHeaderTimeout <= 0 {
		opts.ReadHeaderTimeout = defaultReadHeaderTimeout
	}

	if opts.ReadTimeout <= 0 {
		opts.ReadTimeout = defaultReadTimeout
	}

	if opts.WriteTimeout <= 0 {
		opts.WriteTimeout = defaultWriteTimeout
	}

	if opts.IdleTimeout <= 0 {
		opts.IdleTimeout = defaultIdleTimeout
	}

	return &Server{
		http: &http.Server{
			Addr:              opts.Addr,
			Handler:           newMux(),
			ReadHeaderTimeout: opts.ReadHeaderTimeout,
			ReadTimeout:       opts.ReadTimeout,
			WriteTimeout:      opts.WriteTimeout,
			IdleTimeout:       opts.IdleTimeout,
		},
		log:             opts.Logger,
		shutdownTimeout: opts.ShutdownTimeout,
		started:         make(chan struct{}),
	}
}

// Addr returns the bound address after Run has started listening.
func (s *Server) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.bound
}

// Started is closed once the listener is up.
func (s *Server) Started() <-chan struct{} {
	return s.started
}

// Run listens until ctx is canceled, then shuts down in-flight requests.
func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.bound = ln.Addr().String()
	s.mu.Unlock()
	close(s.started)

	s.log.Info("listening", zap.String("addr", s.bound))

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.http.Serve(ln)
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	case <-ctx.Done():
	}

	s.log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return err
	}

	err = <-errCh
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
