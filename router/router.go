package router

import (
	"ShopManageSystem/handlers"
	"ShopManageSystem/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter(c *gin.Engine) {
	v1 := c.Group("/v1/api")
	{
		pub := v1.Group("/public")
		{
			pub.GET("/getCaptcha", handlers.GetCaptchaCode)
			pub.POST("/register", handlers.Register)
			pub.POST("/loginByPass", handlers.LoginByPass)
			pub.POST("/loginByVerfy", handlers.LoginByVerfiyCode)
			pub.GET("/sendVerfiyCode", handlers.SendVerifyCode)

		}

		auth := v1.Group("auth")
		auth.Use(middlewares.TokenVerify())
		{
			auth.GET("/logout", handlers.Logout)
		}
	}
}
