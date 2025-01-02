package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/pro"
)

type ProjectRepo interface {
	FindProjectByMemId(ctx context.Context, memId int64, page int64, size int64) ([]*pro.ProjectAndMember, int64, error)
}
