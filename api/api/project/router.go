package project

import (
	"github.com/Tuanzi-bug/SyncHub/api/api/midd"
	"github.com/Tuanzi-bug/SyncHub/api/router"
	"github.com/gin-gonic/gin"
	"log"
)

type RouterProject struct {
}

func init() {
	log.Println("init project router")
	ru := &RouterProject{}
	router.Register(ru)
}

func (*RouterProject) Route(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcProjectClient()
	h := New()
	group := r.Group("/project/index")
	group.Use(midd.TokenVerify())
	group.POST("", h.index)
	group1 := r.Group("/project/project")
	group1.Use(midd.TokenVerify())
	group1.POST("/selfList", h.myProjectList)
}
