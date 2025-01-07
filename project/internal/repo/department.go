package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
)

type DepartmentRepo interface {
	FindDepartmentById(ctx context.Context, id int64) (*domain.Department, error)
	FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*domain.Department, error)
	Save(dpm *domain.Department) error
	ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*domain.Department, total int64, err error)
}
