package main

import (
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/user/config"
	"github.com/Tuanzi-bug/SyncHub/user/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	gc := router.RegisterGrpc()
	stop := func() {
		gc.Stop()
	}
	router.RegisterEtcdServer()
	common.Run(r, config.AppConfig.ServerConfig.Name, config.AppConfig.ServerConfig.Addr, stop)
}
