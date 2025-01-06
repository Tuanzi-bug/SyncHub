package dao

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type SourceLinkDao struct {
	conn *gorms.GormConn
}

func (s *SourceLinkDao) Save(ctx context.Context, link *domain.SourceLink) error {
	return s.conn.Session(ctx).Save(&link).Error
}

func (s *SourceLinkDao) FindByTaskCode(ctx context.Context, taskCode int64) (list []*domain.SourceLink, err error) {
	session := s.conn.Session(ctx)
	err = session.Model(&domain.SourceLink{}).Where("link_type=? and link_code=?", "task", taskCode).Find(&list).Error
	return
}

func NewSourceLinkDao() *SourceLinkDao {
	return &SourceLinkDao{
		conn: gorms.New(),
	}
}
