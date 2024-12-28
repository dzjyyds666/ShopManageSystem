package middlewares

import (
	"ShopManageSystem/utils/log/logx"
	"ShopManageSystem/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Recovery 全局异常处理
func Recovery(r *gin.Engine) {
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logx.GetLogger("ShopManage").Errorf("SystemError|%v", recovered)
		c.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.SystemError, "系统异常,请稍后再试", recovered))
	}))
	logx.GetLogger("ShopManage").Infof("Middleware|Recovery|SUCC")
}
