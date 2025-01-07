package dao

import (
	"context"
	"errors"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/gorms"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
	"gorm.io/gorm"
)

type DepartmentDao struct {
	conn *gorms.GormConn
}

func (d DepartmentDao) FindDepartmentById(ctx context.Context, id int64) (dt *domain.Department, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("id=?", id).Find(&dt).Error
	return
}

func (d DepartmentDao) FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*domain.Department, error) {
	session := d.conn.Session(ctx)
	session = session.Model(&domain.Department{}).Where("organization_code=? AND name=?", organizationCode, name)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	var dp *domain.Department
	err := session.Limit(1).Take(&dp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return dp, err
}

func (d DepartmentDao) Save(dpm *domain.Department) error {
	err := d.conn.Session(context.Background()).Save(&dpm).Error
	return err
}

func (d DepartmentDao) ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*domain.Department, total int64, err error) {
	session := d.conn.Session(context.Background())
	session = session.Model(&domain.Department{})
	session = session.Where("organization_code=?", organizationCode)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	err = session.Count(&total).Error
	err = session.Limit(int(size)).Offset(int((page - 1) * size)).Find(&list).Error
	return
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{
		conn: gorms.New(),
	}
}
