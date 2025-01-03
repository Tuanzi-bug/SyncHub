package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/task"
)

type TaskStagesTemplateRepo interface {
	FindInProTemIds(ctx context.Context, ids []int) ([]task.MsTaskStagesTemplate, error)
}
