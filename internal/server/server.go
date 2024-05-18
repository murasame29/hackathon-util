package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/murasame29/hackathon-util/pkg/logger"
	"golang.org/x/sync/errgroup"
)

const (
	DefaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	// srv server
	srv *http.Server
	// shutdown timeout
	shutdownTimeout time.Duration
}

// New はHTTPサーバを生成する
func New(addr string, handler http.Handler, opts ...Option) *Server {
	s := &Server{
		srv: &http.Server{
			Addr:    ":8080",
			Handler: handler,
		},
		shutdownTimeout: DefaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Run はHTTPサーバを起動する。
func (s *Server) Run(ctx context.Context) error {
	logger.Info(ctx, "server listening at ...", logger.Field("addr", s.srv.Addr))
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown はhttp serverを停止する
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info(ctx, "server shutting down ...")
	return s.srv.Shutdown(ctx)
}

// RunWithGraceful はサーバの起動とInterrupt,SIGTERMによる停止信号に対してのGracefulShutdownを提供する
func (s *Server) RunWithGraceful(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return s.Run(ctx)
	})

	group.Go(func() error {
		<-gCtx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	})

	if err := group.Wait(); err != nil && err != context.Canceled {
		logger.Info(ctx, "server shutdown failed, Error: ", logger.Field("err", err))
		os.Exit(1)
	}

	logger.Info(ctx, "server shutdown successfully")
}
