package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/pro"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/task"
)

type TaskStagesTemplateRepo interface {
	FindInProTemIds(ctx context.Context, ids []int) ([]pro.MsTaskStagesTemplate, error)
	FindByProjectTemplateId(ctx context.Context, projectTemplateCode int) (list []*pro.MsTaskStagesTemplate, err error)
}
type TaskStagesRepo interface {
	SaveTaskStages(ctx context.Context, conn database.DbConn, ts *task.TaskStages) error
	FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, pageSize int64) (list []*task.TaskStages, total int64, err error)
	FindById(ctx context.Context, id int) (ts *task.TaskStages, err error)
}

type TaskRepo interface {
	FindTaskByStageCode(ctx context.Context, stageCode int) (list []*pro.Task, err error)
	FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (task *pro.TaskMember, err error)
	FindTaskMaxIdNum(ctx context.Context, projectCode int64) (v *int, err error)
	FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (v *int, err error)
	SaveTask(ctx context.Context, conn database.DbConn, ts *pro.Task) error
	SaveTaskMember(ctx context.Context, conn database.DbConn, tm *pro.TaskMember) error
	FindTaskById(ctx context.Context, taskCode int64) (ts *pro.Task, err error)
	UpdateTaskSort(ctx context.Context, conn database.DbConn, ts *pro.Task) error
	FindTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *pro.Task, err error)
	FindTaskByAssignTo(ctx context.Context, memberId int64, done int, page int64, size int64) ([]*pro.Task, int64, error)
	FindTaskByMemberCode(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*pro.Task, total int64, err error)
	FindTaskByCreateBy(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*pro.Task, total int64, err error)
}
