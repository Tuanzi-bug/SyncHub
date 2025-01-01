package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/user/internal/database"
	"github.com/Tuanzi-bug/SyncHub/user/internal/domain/organization"
)

type OrganizationRepo interface {
	SaveOrganization(conn database.DbConn, ctx context.Context, org *organization.Organization) error
	FindOrganizationByMemId(ctx context.Context, memId int64) ([]*organization.Organization, error)
}
