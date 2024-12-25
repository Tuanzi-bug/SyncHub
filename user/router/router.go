package router

import (
	"github.com/Tuanzi-bug/SyncHub/user/internal/service"
	"github.com/gin-gonic/gin"
)

type Router interface {
	Register(r *gin.Engine)
}
type RegisterRouter struct {
}

func New() RegisterRouter {
	return RegisterRouter{}
}
func (RegisterRouter) Route(router Router, r *gin.Engine) {
	router.Register(r)
}

func InitRouter(r *gin.Engine) {
	router := New()
	//以后的模块路由在这进行注册
	router.Route(service.NewUserHandler(), r)
}
