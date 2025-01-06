package dao

import (
	"context"
	"errors"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/pro"
	"gorm.io/gorm"
)

type TaskDao struct {
	conn *gorms.GormConn
}

func (t TaskDao) FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (v *int, err error) {
	session := t.conn.Session(ctx)
	//select * from
	err = session.Model(&pro.Task{}).
		Where("project_code=? and stage_code=?", projectCode, stageCode).
		Select("max(sort)").Scan(&v).Error
	return
}

func (t TaskDao) SaveTask(ctx context.Context, conn database.DbConn, ts *pro.Task) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&ts).Error
	return err
}

func (t TaskDao) SaveTaskMember(ctx context.Context, conn database.DbConn, tm *pro.TaskMember) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&tm).Error
	return err
}

func (t TaskDao) FindTaskById(ctx context.Context, taskCode int64) (ts *pro.Task, err error) {
	session := t.conn.Session(ctx)
	err = session.Where("id=?", taskCode).Find(&ts).Error
	return
}

func (t TaskDao) UpdateTaskSort(ctx context.Context, conn database.DbConn, ts *pro.Task) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).
		Where("id=?", ts.Id).
		Select("sort", "stage_code").
		Updates(&ts).
		Error
	return err
}

func (t TaskDao) FindTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *pro.Task, err error) {
	session := t.conn.Session(ctx)
	err = session.Where("stage_code=? and sort < ?", stageCode, sort).Order("sort desc").Limit(1).Find(&ts).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (t TaskDao) FindTaskByAssignTo(ctx context.Context, memberId int64, done int, page int64, size int64) (tsList []*pro.Task, total int64, err error) {
	session := t.conn.Session(ctx)
	offset := (page - 1) * size
	err = session.Model(&pro.Task{}).Where("assign_to=? and deleted=0 and done=?", memberId, done).Limit(int(size)).Offset(int(offset)).Find(&tsList).Error
	err = session.Model(&pro.Task{}).Where("assign_to=? and deleted=0 and done=?", memberId, done).Count(&total).Error
	return
}

func (t TaskDao) FindTaskByMemberCode(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*pro.Task, total int64, err error) {
	session := t.conn.Session(ctx)
	offset := (page - 1) * size
	sql := "select a.* from ms_task a,ms_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=? limit ?,?"
	raw := session.Model(&pro.Task{}).Raw(sql, memberId, done, offset, size)
	err = raw.Scan(&tList).Error
	if err != nil {
		return nil, 0, err
	}
	sqlCount := "select count(*) from ms_task a,ms_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=?"
	rawCount := session.Model(&pro.Task{}).Raw(sqlCount, memberId, done)
	err = rawCount.Scan(&total).Error
	return
}

func (t TaskDao) FindTaskByCreateBy(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*pro.Task, total int64, err error) {
	session := t.conn.Session(ctx)
	offset := (page - 1) * size
	err = session.Model(&pro.Task{}).Where("create_by=? and deleted=0 and done=?", memberId, done).Limit(int(size)).Offset(int(offset)).Find(&tList).Error
	err = session.Model(&pro.Task{}).Where("create_by=? and deleted=0 and done=?", memberId, done).Count(&total).Error
	return
}

func (t TaskDao) FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (task *pro.TaskMember, err error) {
	err = t.conn.Session(ctx).
		Where("task_code=? and member_code=?", taskCode, memberId).
		Limit(1).
		Find(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (t TaskDao) FindTaskMaxIdNum(ctx context.Context, projectCode int64) (v *int, err error) {
	session := t.conn.Session(ctx)
	//select * from
	err = session.Model(&pro.Task{}).
		Where("project_code=?", projectCode).
		Select("max(id_num)").Scan(&v).Error
	return
}

func (t TaskDao) FindTaskByStageCode(ctx context.Context, stageCode int) (list []*pro.Task, err error) {
	//select * from ms_task where stage_code=77 and deleted =0 order by sort asc
	session := t.conn.Session(ctx)
	err = session.Model(&pro.Task{}).
		Where("stage_code=? and deleted =0", stageCode).
		Order("sort asc").
		Find(&list).Error
	return
}

func NewTaskDao() *TaskDao {
	return &TaskDao{
		conn: gorms.New(),
	}
}
