package tran

import "github.com/Tuanzi-bug/SyncHub/project/internal/database"

type Transaction interface {
	Action(func(conn database.DbConn) error) error
}
