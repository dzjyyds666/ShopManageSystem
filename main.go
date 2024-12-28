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

// main 程序入口

// @title 超市管理系统
// @version 0.1.0
// @description 超市管瘤系统，供客户端调用

// @contact.name Aaron
// @contact.email duaaron519@gmail.com

// @host http://localhost:8888
// @BasePath /api/v1
func main() {
	c := gin.Default()

	middlewares.Recovery(c)

	router.InitRouter(c)

	err := c.Run(fmt.Sprintf(":%d", config.GlobalConfig.ServerPort))
	if err != nil {
		logx.GetLogger("SH").Errorf("BootStrop|StartError|%v", err)
	}
}

//package main
//
//import "fmt"
//
//func main() {
//
//	var good goods
//
//	good = fruit{
//		name: "fruit",
//	}
//
//	good.getPrice(10, 20, 30)
//
//	good = vegetable{
//		name: "蔬菜",
//	}
//
//	good.getPrice(10, 20, 30)
//
//}
//
//type goods interface {
//	getPrice(prices ...int) int
//}
//
//type fruit struct {
//	name  string
//	price int
//}
//
//func (fruit) getPrice(prices ...int) int {
//	var price int
//	for _, item := range prices {
//		price = price + item
//	}
//	fmt.Printf("水果的价格是 %v\n", price)
//	return price
//}
//
//type vegetable struct {
//	name  string
//	price int
//}
//
//func (vegetable) getPrice(prices ...int) int {
//	var price int
//	for _, item := range prices {
//		price = price + item
//	}
//	fmt.Printf("蔬菜的价格是 %v\n", price)
//	return price
//}
