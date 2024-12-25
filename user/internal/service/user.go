package service

import (
	common "github.com/Tuanzi-bug/SyncHub/common"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Signup()
}

type userService struct {
}

func (u *userService) GetCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	ctx.JSON(200, result.Success("123456"))
}
