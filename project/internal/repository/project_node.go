package repository

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/project/internal/dao"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain"
	"github.com/Tuanzi-bug/SyncHub/project/internal/repo"
	"github.com/Tuanzi-bug/SyncHub/project/pkg/grpc_errs"
)

type ProjectNodeDomain struct {
	projectNodeRepo repo.ProjectNodeRepo
}

func (d *ProjectNodeDomain) TreeList() ([]*domain.ProjectNodeTree, *errs.BError) {
	//node表都查出来 转换成treelist结构
	list, err := d.projectNodeRepo.FindAll(context.Background())
	if err != nil {
		return nil, grpc_errs.DBError
	}
	return domain.ToNodeTreeList(list), nil
}
func (d *ProjectNodeDomain) NodeList() ([]*domain.ProjectNode, *errs.BError) {
	list, err := d.projectNodeRepo.FindAll(context.Background())
	if err != nil {
		return nil, grpc_errs.DBError
	}
	return list, nil
}

func NewProjectNodeDomain() *ProjectNodeDomain {
	return &ProjectNodeDomain{
		projectNodeRepo: dao.NewProjectNodeDao(),
	}
}
