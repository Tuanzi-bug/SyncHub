package main

import (
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	common.Run(r, "project-user", ":9080")
}
