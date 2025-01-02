package user

import (
	"github.com/Tuanzi-bug/SyncHub/api/api/midd"
	"github.com/Tuanzi-bug/SyncHub/api/api/user_rpc"
	"github.com/Tuanzi-bug/SyncHub/api/router"
	"github.com/gin-gonic/gin"
	"log"
)

type RouterUser struct {
}

func init() {
	log.Println("init user_proto router")
	ru := &RouterUser{}
	router.Register(ru)
}

func (*RouterUser) Route(r *gin.Engine) {
	//初始化grpc的客户端连接
	user_rpc.InitRpcUserClient()
	h := New()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.Register)
	r.POST("/project/login", h.login)
	org := r.Group("/project/organization")
	org.Use(midd.TokenVerify())
	org.POST("/_getOrgList", h.myOrgList)
}
