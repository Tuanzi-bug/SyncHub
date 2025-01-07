package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type ProjectAuthRepo interface {
	FindAuthList(ctx context.Context, orgCode int64) (list []*domain.ProjectAuth, err error)
	FindAuthListPage(ctx context.Context, orgCode int64, page int64, pageSize int64) (list []*domain.ProjectAuth, total int64, err error)
}
