package project

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/api/pkg/params"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/grpc/project"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type HandlerProject struct {
}

func New() *HandlerProject {
	return &HandlerProject{}
}

func (p *HandlerProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}
	indexResponse, err := ProjectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(indexResponse.Menus))
}
func (p *HandlerProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	//1. 获取参数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberIdStr, _ := c.Get("memberId")
	memberId := memberIdStr.(int64)
	page := &params.Page{}
	page.Bind(c)
	msg := &project.ProjectRpcMessage{MemberId: memberId, Page: page.Page, PageSize: page.PageSize}
	myProjectResponse, err := ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	if myProjectResponse.Pm == nil {
		myProjectResponse.Pm = []*project.ProjectMessage{}
	}
	var pms []*params.ProjectAndMember
	copier.Copy(&pms, myProjectResponse.Pm)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": myProjectResponse.Total,
	}))
}
