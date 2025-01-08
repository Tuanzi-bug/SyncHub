package menu_service_v1

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/grpc/menu"
	"github.com/Tuanzi-bug/SyncHub/project/internal/dao"
	"github.com/Tuanzi-bug/SyncHub/project/internal/database/tran"
	"github.com/Tuanzi-bug/SyncHub/project/internal/repo"
	"github.com/Tuanzi-bug/SyncHub/project/internal/repository"
	"github.com/jinzhu/copier"
)

type MenuService struct {
	menu.UnimplementedMenuServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	menuDomain  *repository.MenuDomain
}

func New() *MenuService {
	return &MenuService{
		cache:       dao.Rc,
		transaction: dao.NewTransaction(),
		menuDomain:  repository.NewMenuDomain(),
	}
}
func (d *MenuService) MenuList(ctx context.Context, msg *menu.MenuReqMessage) (*menu.MenuResponseMessage, error) {
	list, err := d.menuDomain.MenuTreeList()
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var mList []*menu.MenuMessage
	copier.Copy(&mList, list)
	return &menu.MenuResponseMessage{List: mList}, nil
}
