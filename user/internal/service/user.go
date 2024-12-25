package service

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/user/internal/dao"
	"github.com/Tuanzi-bug/SyncHub/user/internal/errs"
	"github.com/Tuanzi-bug/SyncHub/user/internal/repo"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type UserHandler struct {
	cache repo.Cache
}

func NewUserHandler() *UserHandler {
	return &UserHandler{cache: dao.Rc}
}

func (h *UserHandler) Register(r *gin.Engine) {
	g := r.Group("/project/login")
	g.POST("/getCaptcha", h.GetCaptcha)
}

func (h *UserHandler) GetCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	//1. 获取参数
	mobile := ctx.PostForm("mobile")
	//2. 验证手机合法性
	if !common.VerifyMobile(mobile) {
		ctx.JSON(200, result.Fail(errs.LoginMobileNotLegal, "不合法"))
		return
	}
	//3.生成验证码
	code := "123456"
	//4. 发送验证码
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("调用短信平台发送短信")
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := h.cache.Put(c, mobile, code, 5*time.Minute)
		if err != nil {
			log.Println("验证码存储失败")
		}
	}()
	ctx.JSON(200, result.Success("123456"))
}
