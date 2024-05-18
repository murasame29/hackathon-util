package server

import (
	"context"
)

func (s *Server) OpenBot(ctx context.Context) error {
	return s.ss.Open(ctx)
}
