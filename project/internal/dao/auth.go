package dao

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type ProjectAuthDao struct {
	conn *gorms.GormConn
}

func (p ProjectAuthDao) FindAuthList(ctx context.Context, orgCode int64) (list []*domain.ProjectAuth, err error) {
	session := p.conn.Session(ctx)
	err = session.Model(&domain.ProjectAuth{}).Where("organization_code=?", orgCode).Find(&list).Error
	return
}

func (p ProjectAuthDao) FindAuthListPage(ctx context.Context, orgCode int64, page int64, pageSize int64) (list []*domain.ProjectAuth, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.Model(&domain.ProjectAuth{}).
		Where("organization_code=?", orgCode).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	err = session.Model(&domain.ProjectAuth{}).
		Where("organization_code=?", orgCode).
		Count(&total).Error
	return
}

func NewProjectAuthDao() *ProjectAuthDao {
	return &ProjectAuthDao{
		conn: gorms.New(),
	}
}
