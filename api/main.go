package main

import (
	_ "github.com/Tuanzi-bug/SyncHub/api/api"
	"github.com/Tuanzi-bug/SyncHub/api/router"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/user/config"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//路由
	router.InitRouter(r)
	common.Run(r, config.AppConfig.ServerConfig.Name, config.AppConfig.ServerConfig.Addr, nil)
}
