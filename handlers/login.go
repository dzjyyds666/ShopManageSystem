package handlers

import (
	"ShopManageSystem/config"
	"ShopManageSystem/database"
	"ShopManageSystem/models"
	"ShopManageSystem/utils/jwt"
	"ShopManageSystem/utils/log/logx"
	"ShopManageSystem/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type loginInfo struct {
	Email       string `json:"email"`        // 邮箱
	Password    string `json:"password"`     // 密码
	VerifyCode  string `json:"verify_code"`  // 邮箱验证码信息
	CaptchaId   string `json:"captcha_id"`   // 验证码id
	CaptchaCode string `json:"captcha_code"` // 验证码
}

// LoginByPass
// @Summary 用户登录
// @Tags 登录
// @Accept json
// @Produce json
// @Param loginInfo body handlers.loginInfo true "登录信息"
// @Success 200 {object} response.Result "登录成功"
// @Router /login [post]
func LoginByPass(ctx *gin.Context) {
	var logininfo loginInfo

	err := ctx.ShouldBindJSON(&logininfo)
	if err != nil {
		logx.GetLogger("SH").Errorf("Handler|LoginByPass|ParamsError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.ParamError, "参数错误", nil))
		return
	}

	// 验证验证码
	getCaptchaCode := database.RDB[0].Get(ctx, fmt.Sprintf(database.Redis_Captcha_Key, logininfo.CaptchaId))

	if getCaptchaCode.Val() != logininfo.CaptchaCode {
		logx.GetLogger("SH").Errorf("Handler|LoginByPass|CaptchaError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.ParamError, "验证码错误", nil))
		return
	}

	var userInfo models.UserInfo
	database.MyDB.Where("email = ?", logininfo.Email).First(&userInfo)

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(logininfo.Password))
	if err != nil {
		logx.GetLogger("SH").Errorf("Handler|LoginByPass|PasswordError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.PasswordError, "密码错误", nil))
		return
	}

	// 把用户token存入redis
	token := jwt.NewJWTUtils().CreateJWT(userInfo.UserId)
	tokenExpirationTime := time.Duration(config.GlobalConfig.JWT.ExpirationTime) * time.Hour
	err = database.RDB[0].Set(ctx, fmt.Sprintf(database.Redis_Token_Key, userInfo.UserId), token, tokenExpirationTime).Err()
	if err != nil {
		logx.GetLogger("SH").Errorf("Handler|LoginByVerfiyCode|RedisSetError|%v", err)
		panic(err)
	}

	// 返回给用户数据
	ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.RequestSuccess, "登录成功", models.UserInfo{
		UserId:   userInfo.UserId,
		Email:    userInfo.Email,
		UserName: userInfo.UserName,
		Avatar:   userInfo.Avatar,
	}))
}

func LoginByVerfiyCode(ctx *gin.Context) {
	var logininfo loginInfo
	err := ctx.ShouldBindJSON(&logininfo)
	if err != nil {
		logx.GetLogger("SH").Errorf("Handler|LoginByVerfiyCode|ParamsError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.ParamError, "参数错误", nil))
		return
	}

	// 验证验证码
	getCaptchaCode := database.RDB[0].Get(ctx, fmt.Sprintf(database.Redis_Captcha_Key, logininfo.CaptchaId))

	if getCaptchaCode.Val() != logininfo.CaptchaCode {
		logx.GetLogger("SH").Errorf("Handler|LoginByPass|CaptchaError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.ParamError, "图片验证码错误", nil))
		return
	}

	// 验证邮箱验证码
	result := database.RDB[0].Get(ctx, fmt.Sprintf(database.Redis_Verification_Code_Key, logininfo.Email))
	if result.Val() != logininfo.VerifyCode {
		logx.GetLogger("SH").Errorf("Handler|LoginByVerfiyCode|VerifyCodeError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.ParamError, "邮箱验证码错误", nil))
		return
	}

	// 获取用户数据
	var userInfo models.UserInfo
	err = database.MyDB.Where("email = ?", logininfo.Email).First(&userInfo).Error
	if err != nil {
		logx.GetLogger("SH").Errorf("Handler|LoginByVerfiyCode|GetUserInfoError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.ParamError, "用户不存在", nil))
		return
	}

	// 把用户token存入redis
	token := jwt.NewJWTUtils().CreateJWT(userInfo.UserId)
	tokenExpirationTime := time.Duration(config.GlobalConfig.JWT.ExpirationTime) * time.Hour
	err = database.RDB[0].Set(ctx, fmt.Sprintf(database.Redis_Token_Key, userInfo.UserId), token, tokenExpirationTime).Err()
	if err != nil {
		logx.GetLogger("SH").Errorf("Handler|LoginByVerfiyCode|RedisSetError|%v", err)
		panic(err)
	}

	// 返回给用户数据
	ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.RequestSuccess, "登录成功", models.UserInfo{
		UserId:   userInfo.UserId,
		Email:    userInfo.Email,
		UserName: userInfo.UserName,
		Avatar:   userInfo.Avatar,
	}))
}

func Logout(ctx *gin.Context) {

	// 删除redis的token信息
	err := database.RDB[0].Del(ctx, fmt.Sprintf(database.Redis_Token_Key, ctx.GetString("user_id"))).Err()
	if err != nil {
		logx.GetLogger("SH").Errorf("Handler|Logout|RedisDelError|%v", err)
		ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.SystemError, "系统异常", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewResult(response.EnmuHttptatus.RequestSuccess, "退出成功", nil))
}
