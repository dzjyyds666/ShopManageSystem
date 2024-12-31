package router

import (
	"ShopManageSystem/handlers"
	"ShopManageSystem/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(c *gin.Engine) {
	c.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := c.Group("/api/v1")
	{
		pub := v1.Group("")
		{
			pub.GET("/getCaptcha", handlers.GetCaptchaCode)
			pub.POST("/register", handlers.Register)
			pub.POST("/loginByPass", handlers.LoginByPass)
			pub.POST("/loginByVerfy", handlers.LoginByVerfiyCode)
			pub.GET("/sendVerfiyCode", handlers.SendVerifyCode)
			pub.POST("/upload/file", handlers.UploadFile)
		}

		auth := v1.Group("")
		// todo 装饰器模式
		auth.Use(middlewares.TokenVerify()) // 添加验证token中间件
		// 动态添加一些别的功能，比如接口线流，代理
		{
			auth.GET("/logout", handlers.Logout)
			auth.GET("/goods/info/:goods_id", handlers.GetGoodsInfo)
			auth.GET("/goods/list", handlers.GetGoodsListByPage)
			auth.POST("/goods/onShelves", handlers.GoodsOnShelves)

			auth.GET("/getUserInfo", handlers.GetUserInfo)
			auth.GET("/getUserList", handlers.GetUserList)
			auth.POST("/updateUserInfo", handlers.UpdateUserInfo)
			auth.GET("/changeRole", handlers.ChangeUserRole)

			auth.POST("/goods/initPayment", handlers.InitBuyGoods)
			auth.GET("/goods/payForOrder", handlers.PayForOrder)
			auth.GET("/goods/hetPayTotal", handlers.GetPayTotal)
			auth.GET("/goods/completeOrder/:order_id", handlers.CompleteOrder)

			auth.GET("/goods/markDiscountGoods", handlers.MarkDiscountGoods)
			auth.GET("/goods/cancelDiscountGoods", handlers.CancelDiscountGoods)
			auth.GET("/goods/discount/:discount", handlers.Discount)

			auth.POST("/type/create", handlers.CreateType)
			auth.GET("/type/list", handlers.GetTypeList)
			auth.GET("/type/delete/:type_id", handlers.DeleteType)
			auth.GET("/type/search", handlers.SearchType)
			auth.GET("/type/getinfo/:type_id", handlers.GetTypeInfo)
			auth.GET("/type/getDeleteTypeList", handlers.GetDeleteTypeList)
		}
	}
}
