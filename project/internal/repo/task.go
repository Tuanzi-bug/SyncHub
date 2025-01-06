package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type TaskStagesTemplateRepo interface {
	FindInProTemIds(ctx context.Context, ids []int) ([]domain.MsTaskStagesTemplate, error)
	FindByProjectTemplateId(ctx context.Context, projectTemplateCode int) (list []*domain.MsTaskStagesTemplate, err error)
}
type TaskStagesRepo interface {
	SaveTaskStages(ctx context.Context, conn database.DbConn, ts *domain.TaskStages) error
	FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, pageSize int64) (list []*domain.TaskStages, total int64, err error)
	FindById(ctx context.Context, id int) (ts *domain.TaskStages, err error)
}

type TaskRepo interface {
	FindTaskByStageCode(ctx context.Context, stageCode int) (list []*domain.Task, err error)
	FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (task *domain.TaskMember, err error)
	FindTaskMaxIdNum(ctx context.Context, projectCode int64) (v *int, err error)
	FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (v *int, err error)
	SaveTask(ctx context.Context, conn database.DbConn, ts *domain.Task) error
	SaveTaskMember(ctx context.Context, conn database.DbConn, tm *domain.TaskMember) error
	FindTaskById(ctx context.Context, taskCode int64) (ts *domain.Task, err error)
	UpdateTaskSort(ctx context.Context, conn database.DbConn, ts *domain.Task) error
	FindTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *domain.Task, err error)
	FindTaskByAssignTo(ctx context.Context, memberId int64, done int, page int64, size int64) ([]*domain.Task, int64, error)
	FindTaskByMemberCode(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*domain.Task, total int64, err error)
	FindTaskByCreateBy(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*domain.Task, total int64, err error)
	FindTaskMemberPage(ctx context.Context, taskCode int64, page int64, size int64) (list []*domain.TaskMember, total int64, err error)
	FindTaskByIds(background context.Context, taskIdList []int64) (list []*domain.Task, err error)
}
