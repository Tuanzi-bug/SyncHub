package project

import (
	"context"
	"encoding/json"
	"github.com/Tuanzi-bug/SyncHub/api/pkg/params"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/grpc/auth"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type HandlerAuth struct {
}

func NewAuth() *HandlerAuth {
	return &HandlerAuth{}
}

func (a *HandlerAuth) authList(c *gin.Context) {
	result := &common.Result{}
	organizationCode := c.GetString("organizationCode")
	var page = &params.Page{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		OrganizationCode: organizationCode,
		Page:             page.Page,
		PageSize:         page.PageSize,
	}

	//1. 去auth表查询authList
	response, err := AuthServiceClient.AuthList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	// 处理返回结果
	var authList []*params.ProjectAuth
	copier.Copy(&authList, response.List)
	if authList == nil {
		authList = []*params.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": response.Total,
		"list":  authList,
		"page":  page.Page,
	}))
}

func (a *HandlerAuth) apply(c *gin.Context) {
	result := &common.Result{}
	var req *params.ProjectAuthReq
	c.ShouldBind(&req)
	var nodes []string
	if req.Nodes != "" {
		json.Unmarshal([]byte(req.Nodes), &nodes)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	msg := &auth.AuthReqMessage{
		Action: req.Action,
		AuthId: req.Id,
		Nodes:  nodes,
	}
	applyResponse, err := AuthServiceClient.Apply(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var list []*params.ProjectNodeAuthTree
	copier.Copy(&list, applyResponse.List)
	var checkedList []string
	copier.Copy(&checkedList, applyResponse.CheckedList)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":        list,
		"checkedList": checkedList,
	}))
}

func (a *HandlerAuth) GetAuthNodes(c *gin.Context) ([]string, error) {
	memberId := c.GetInt64("memberId")
	msg := &auth.AuthReqMessage{
		MemberId: memberId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := AuthServiceClient.AuthNodesByMemberId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		return nil, errs.NewError(errs.ErrorCode(code), msg)
	}
	return response.List, err
}
