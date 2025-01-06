package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type FileRepo interface {
	Save(ctx context.Context, file *domain.File) error
	FindByIds(background context.Context, ids []int64) (list []*domain.File, err error)
}
