package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type ProjectLogRepo interface {
	FindLogByTaskCode(ctx context.Context, taskCode int64, comment int) (list []*domain.ProjectLog, total int64, err error)
	FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) (list []*domain.ProjectLog, total int64, err error)
	SaveProjectLog(pl *domain.ProjectLog)
	FindLogByMemberCode(background context.Context, memberId int64, page int64, size int64) (list []*domain.ProjectLog, total int64, err error)
}
