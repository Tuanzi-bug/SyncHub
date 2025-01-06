package dao

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type FileDao struct {
	conn *gorms.GormConn
}

func (f *FileDao) FindByIds(ctx context.Context, ids []int64) (list []*domain.File, err error) {
	session := f.conn.Session(ctx)
	err = session.Model(&domain.File{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func (f *FileDao) Save(ctx context.Context, file *domain.File) error {
	err := f.conn.Session(ctx).Save(&file).Error
	return err
}

func NewFileDao() *FileDao {
	return &FileDao{
		conn: gorms.New(),
	}
}
