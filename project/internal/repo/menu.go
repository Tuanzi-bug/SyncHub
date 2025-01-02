package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/project/internal/domain/menu"
)

type MenuRepo interface {
	FindMenus(ctx context.Context) ([]*menu.ProjectMenu, error)
}
