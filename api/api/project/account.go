package project

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/api/pkg/params"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/grpc/account"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type HandlerAccount struct {
}

func NewAccount() *HandlerAccount {
	return &HandlerAccount{}
}

func (a *HandlerAccount) account(c *gin.Context) {
	//接收请求参数  一些参数的校验 可以放在api这里
	result := &common.Result{}
	var req *params.AccountReq
	_ = c.ShouldBind(&req)
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	//调用project模块 查询账户列表
	msg := &account.AccountReqMessage{
		MemberId:         memberId,
		OrganizationCode: c.GetString("organizationCode"),
		Page:             int64(req.Page),
		PageSize:         int64(req.PageSize),
		SearchType:       int32(req.SearchType),
		DepartmentCode:   req.DepartmentCode,
	}
	response, err := AccountServiceClient.Account(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//返回数据
	var list []*params.MemberAccount
	copier.Copy(&list, response.AccountList)
	if list == nil {
		list = []*params.MemberAccount{}
	}
	var authList []*params.ProjectAuth
	copier.Copy(&authList, response.AuthList)
	if authList == nil {
		authList = []*params.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total":    response.Total,
		"page":     req.Page,
		"list":     list,
		"authList": authList,
	}))
}
