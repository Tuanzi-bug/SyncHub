package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type ProjectRepo interface {
	FindProjectByMemId(ctx context.Context, memId int64, condition string, page int64, size int64) ([]*domain.ProjectAndMember, int64, error)
	FindCollectProjectByMemId(ctx context.Context, memberId int64, page int64, size int64) ([]*domain.ProjectAndMember, int64, error)
	SaveProject(conn database.DbConn, ctx context.Context, pr *domain.Project) error
	SaveProjectMember(conn database.DbConn, ctx context.Context, pm *domain.ProjectMember) error
	FindProjectByPIdAndMemId(ctx context.Context, projectCode int64, memberId int64) (*domain.ProjectAndMember, error)
	FindCollectByPidAndMemId(ctx context.Context, projectCode int64, memberId int64) (bool, error)
	UpdateDeletedProject(ctx context.Context, code int64, deleted bool) error
	SaveProjectCollect(ctx context.Context, pc *domain.ProjectCollection) error
	DeleteProjectCollect(ctx context.Context, memId int64, projectCode int64) error
	UpdateProject(ctx context.Context, proj *domain.Project) error
	FindProjectMemberByPid(ctx context.Context, projectCode int64) (list []*domain.ProjectMember, total int64, err error)
	FindProjectById(ctx context.Context, projectCode int64) (pj *domain.Project, err error)
	FindProjectByIds(ctx context.Context, pids []int64) (list []*domain.Project, err error)
}

type ProjectTemplateRepo interface {
	FindProjectTemplateSystem(ctx context.Context, page int64, size int64) ([]domain.ProjectTemplate, int64, error)
	FindProjectTemplateCustom(ctx context.Context, memId int64, organizationCode int64, page int64, size int64) ([]domain.ProjectTemplate, int64, error)
	FindProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) ([]domain.ProjectTemplate, int64, error)
}
