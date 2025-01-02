package main

import (
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/project/config"
	"github.com/Tuanzi-bug/SyncHub/project/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//路由
	router.InitRouter(r)
	//grpc服务注册
	gc := router.RegisterGrpc()
	//grpc服务注册到etcd
	router.RegisterEtcdServer()

	stop := func() {
		gc.Stop()
	}
	common.Run(r, config.AppConfig.ServerConfig.Name, config.AppConfig.ServerConfig.Addr, stop)
}
