package project

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/api/pkg/params"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/grpc/menu"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type HandlerMenu struct {
}

func NewMenu() *HandlerMenu {
	return &HandlerMenu{}
}

func (m HandlerMenu) menuList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res, err := MenuServiceClient.MenuList(ctx, &menu.MenuReqMessage{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var list []*params.Menu
	copier.Copy(&list, res.List)
	if list == nil {
		list = []*params.Menu{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}
