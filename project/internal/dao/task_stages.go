package dao

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type TaskStagesDao struct {
	conn *gorms.GormConn
}

func (t TaskStagesDao) SaveTaskStages(ctx context.Context, conn database.DbConn, ts *domain.TaskStages) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&ts).Error
	return err
}

func (t TaskStagesDao) FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, pageSize int64) (list []*domain.TaskStages, total int64, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&domain.TaskStages{}).
		Where("project_code=?", projectCode).
		Order("sort asc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	err = session.Model(&domain.TaskStages{}).
		Where("project_code=?", projectCode).
		Count(&total).Error
	return
}

func (t TaskStagesDao) FindById(ctx context.Context, id int) (ts *domain.TaskStages, err error) {
	err = t.conn.Session(ctx).Where("id=?", id).Find(&ts).Error
	return
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{
		conn: gorms.New(),
	}
}
