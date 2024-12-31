package handlers

import (
	"ShopManageSystem/database"
	"ShopManageSystem/models"
	"ShopManageSystem/utils/log/logx"
	"ShopManageSystem/utils/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserInfo(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")

	var userInfo models.UserInfo

	err := database.MyDB.
		Select("user_id", "user_name", "email", "role", "avatar").
		Where("user_id = ?", userId).
		First(&userInfo).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("GetUserInfo|MySqlError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取用户信息失败", nil))
		ctx.Abort()
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取用户信息成功", userInfo))
}

func GetUserList(ctx *gin.Context) {
	var userInfo []models.UserInfo
	err := database.MyDB.Select("user_id", "user_name", "email", "role", "avatar").Find(&userInfo).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("GetUserList|MySqlError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取用户列表失败", nil))
		ctx.Abort()
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取用户列表成功", userInfo))
}

func UpdateUserInfo(ctx *gin.Context) {
	var userInfo models.UserInfo
	err := ctx.ShouldBindJSON(&userInfo)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("UpdateUserInfo|ParamsError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "参数错误", nil))
		ctx.Abort()
	}

	// 对密码进行加密
	password, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	userInfo.Password = string(password)

	err = database.MyDB.Model(&userInfo).Where("user_id = ?", userInfo.UserId).Updates(&userInfo).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("UpdateUserInfo|MySqlError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "更新用户信息失败", nil))
		ctx.Abort()
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "更新用户信息成功", nil))
}

func ChangeUserRole(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	role := ctx.Query("role")

	err := database.MyDB.Model(&models.UserInfo{}).Where("user_id = ?", userId).Update("role", role).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("ChangeUserRole|MySqlError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "更新用户角色失败", nil))
		ctx.Abort()
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "更新用户角色成功", nil))
}
