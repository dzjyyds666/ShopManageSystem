package middlewares

import (
	"ShopManageSystem/database"
	"ShopManageSystem/utils/jwt"
	"ShopManageSystem/utils/log/logx"
	"ShopManageSystem/utils/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// token校验中间键
func TokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if len(token) <= 0 {
			logx.GetLogger("ShopManage").Errorf("TokenVerify|token为空,请重新登录")
			c.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.TokenInvalid, "token为空,请重新登录", nil))
			// 阻止执行
			c.Abort()
			return
		}

		logx.GetLogger("ShopManage").Infof("TokenVerify|token: %v", token)

		jwt := jwt.NewJWTUtils()

		// 解析token
		claims, err := jwt.ParseJWT(token)
		if err != nil {
			c.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.TokenInvalid, "token校验失败,请重新登录", nil))
			// 阻止执行
			c.Abort()
			return
		}
		userId := claims.UserId

		// 从redis中获取token
		redisResult := database.RDB[0].Get(c.Request.Context(), fmt.Sprintf(database.Redis_Token_Key, userId))
		err = redisResult.Err()
		if err != nil {
			if err == redis.Nil {
				c.JSON(http.StatusOK, gin.H{
					"code": response.EnmuHttptatus.TokenExpired,
					"msg":  "token过期,请重新登录",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": response.EnmuHttptatus.RedisError,
					"msg":  "redis错误",
				})
			}
			// 阻止执行
			c.Abort()
			return
		}
		c.Set("user_id", userId)
		c.Next()
	}
}
