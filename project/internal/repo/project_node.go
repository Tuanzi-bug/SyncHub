package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type ProjectNodeRepo interface {
	FindAll(ctx context.Context) (list []*domain.ProjectNode, err error)
}
