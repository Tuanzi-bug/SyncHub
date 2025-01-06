package dao

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/pro"
)

type TaskStagesTemplateDao struct {
	conn *gorms.GormConn
}

func (t *TaskStagesTemplateDao) FindByProjectTemplateId(ctx context.Context, projectTemplateCode int) (list []*pro.MsTaskStagesTemplate, err error) {
	session := t.conn.Session(ctx)
	err = session.
		Model(&pro.MsTaskStagesTemplate{}).
		Where("project_template_code=?", projectTemplateCode).
		Order("sort desc,id asc").
		Find(&list).
		Error
	return
}

func (t *TaskStagesTemplateDao) FindInProTemIds(ctx context.Context, ids []int) ([]pro.MsTaskStagesTemplate, error) {
	var tsts []pro.MsTaskStagesTemplate
	session := t.conn.Session(ctx)
	err := session.Where("project_template_code in ?", ids).Find(&tsts).Error
	return tsts, err
}

func NewTaskStagesTemplateDao() *TaskStagesTemplateDao {
	return &TaskStagesTemplateDao{
		conn: gorms.New(),
	}
}
