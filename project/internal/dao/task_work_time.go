package dao

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type TaskWorkTimeDao struct {
	conn *gorms.GormConn
}

func (t *TaskWorkTimeDao) Save(ctx context.Context, twt *domain.TaskWorkTime) error {
	session := t.conn.Session(ctx)
	err := session.Save(&twt).Error
	return err
}

func (t *TaskWorkTimeDao) FindWorkTimeList(ctx context.Context, taskCode int64) (list []*domain.TaskWorkTime, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&domain.TaskWorkTime{}).Where("task_code=?", taskCode).Find(&list).Error
	return
}

func NewTaskWorkTimeDao() *TaskWorkTimeDao {
	return &TaskWorkTimeDao{
		conn: gorms.New(),
	}
}
