package main

import (
	"ShopManageSystem/config"
	_ "ShopManageSystem/database"
	"ShopManageSystem/middlewares"
	"ShopManageSystem/router"
	"ShopManageSystem/utils/log/logx"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	c := gin.Default()

	middlewares.Recovery(c)

	router.InitRouter(c)

	err := c.Run(fmt.Sprintf(":%d", config.GlobalConfig.ServerPort))
	if err != nil {
		logx.GetLogger("SH").Errorf("BootStrop|StartError|%v", err)
	}
}
