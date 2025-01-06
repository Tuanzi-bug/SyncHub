package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type TaskWorkTimeRepo interface {
	Save(ctx context.Context, twt *domain.TaskWorkTime) error
	FindWorkTimeList(ctx context.Context, taskCode int64) (list []*domain.TaskWorkTime, err error)
}
