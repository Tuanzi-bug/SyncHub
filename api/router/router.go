package router

import (
	"github.com/gin-gonic/gin"
)

// Router 接口
type Router interface {
	Route(r *gin.Engine)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	//rg := New()
	//rg.Route(&user.RouterUser{}, r)
	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}
