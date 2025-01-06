package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type SourceLinkRepo interface {
	Save(ctx context.Context, link *domain.SourceLink) error
	FindByTaskCode(ctx context.Context, taskCode int64) (list []*domain.SourceLink, err error)
}
