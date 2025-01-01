package repo

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/user/internal/database"
	"github.com/Tuanzi-bug/SyncHub/user/internal/domain/member"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, pwd string) (mem *member.Member, err error)
}
