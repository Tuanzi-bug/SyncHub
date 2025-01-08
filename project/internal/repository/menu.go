package repository

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/project/internal/dao"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/menu"
	"github.com/Tuanzi-bug/SyncHub/project/internal/repo"
	"github.com/Tuanzi-bug/SyncHub/project/pkg/grpc_errs"
	"time"
)

type MenuDomain struct {
	menuRepo repo.MenuRepo
}

func (d *MenuDomain) MenuTreeList() ([]*menu.ProjectMenuChild, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	menus, err := d.menuRepo.FindMenus(c)
	if err != nil {
		return nil, grpc_errs.DBError
	}
	menuChildren := menu.CovertChild(menus)
	return menuChildren, nil
}

func NewMenuDomain() *MenuDomain {
	return &MenuDomain{
		menuRepo: dao.NewMenuDao(),
	}
}
