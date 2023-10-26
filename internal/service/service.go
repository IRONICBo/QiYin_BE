package service

import (
	"context"
)

// Service service.
type Service struct {
	ctx context.Context
}

// GetCtx get context.
func (svc *Service) GetCtx() context.Context {
	return svc.ctx
}
