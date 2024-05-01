package server

import (
	"time"
)

type Option func(s *Server)

func WithShutdownTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = t
	}
}
func WithReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.srv.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.srv.WriteTimeout = t
	}
}

func WithOnShutdown(f func()) Option {
	return func(s *Server) {
		s.srv.RegisterOnShutdown(f)
	}
}

// ...必要に応じて追加していく...
