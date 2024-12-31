package handlers

import (
	"ShopManageSystem/database"
	"ShopManageSystem/models"
	"ShopManageSystem/utils/log/logx"
	"ShopManageSystem/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

// Create 创建商品分类
// @Summary 创建商品分类
// @Description 创建商品分类
// @Tags type
// @Accept json
// @Produce json
// @Param goods_type body models.GoodsType true "商品分类"
// @Router /type/create [post]
func CreateType(ctx *gin.Context) {
	var goodsType models.GoodsType
	err := ctx.ShouldBindJSON(&goodsType)
	if err != nil {
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "参数错误", nil))
		ctx.Abort()
	}

	database.RDB[0].Del(ctx, fmt.Sprintf(database.Redis_delete_type_key, goodsType.TypeName)).Err()

	typeIdUUID, _ := uuid.NewUUID()
	typeId := strings.ReplaceAll(typeIdUUID.String(), "-", "")
	goodsType.TypeId = typeId
	err = database.MyDB.Create(&goodsType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "分类已存在", nil))
			ctx.Abort()
		}
		panic(err)
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "创建成功", nil))
}

// GetTypeList 获取商品分类列表
// @Summary 获取商品分类列表
// @Description 获取商品分类列表
// @Tags type
// @Accept json
// @Produce json
// @Router /type/list [get]
func GetTypeList(ctx *gin.Context) {
	var goodsType []models.GoodsType
	err := database.MyDB.Find(&goodsType).Error
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取成功", goodsType))
}

// DeleteType 删除商品分类
// @Summary 删除商品分类
// @Description 删除商品分类
// @Tags type
// @Accept json
// @Produce json
// @Param type_id path string true "type_id"
// @Router /type/delete/{type_id} [get]
// todo 备忘录模式
func DeleteType(ctx *gin.Context) {
	typeId := ctx.Param("type_id")
	var goodsType models.GoodsType
	err := database.MyDB.Where("type_id = ?", typeId).First(&goodsType).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|DeleteType|MySqlError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "该分类不存在", nil))
		ctx.Abort()
	}

	marshal, _ := json.Marshal(&goodsType)

	// 把信息存入redis中做备份
	err = database.RDB[0].Set(ctx, fmt.Sprintf(database.Redis_delete_type_key, goodsType.TypeName), marshal, time.Hour*24*15).Err()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|DeleteType|RedisError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "删除失败，备份删除失败", nil))
		ctx.Abort()
	}

	err = database.MyDB.Delete(&models.GoodsType{}, "type_id=?", typeId).Error
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "删除成功", nil))
}

// GetDeleteTypeList 获取删除的商品分类
func GetDeleteTypeList(ctx *gin.Context) {
	cursor := uint64(0)
	keys := []string{}
	for {
		var keysPart []string
		var err error
		keysPart, cursor, err = database.RDB[0].Scan(ctx, cursor, fmt.Sprintf(database.Redis_delete_type_key, "*"), 0).Result()
		if err != nil {
			logx.GetLogger("ShopManage").Errorf("Handler|DeleteType|RedisError|%v", err)
			ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "查询失败", nil))
			ctx.Abort()
			return
		}
		keys = append(keys, keysPart...)
		if cursor == 0 {
			break
		}
	}
	var deleteType []models.GoodsType
	for _, key := range keys {
		result, err := database.RDB[0].Get(ctx, key).Result()
		if err != nil {
			logx.GetLogger("ShopManage").Errorf("Handler|GetDeleteTypeList|RedisError|%v", err)
			ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "查询失败", nil))
			ctx.Abort()
		}

		var goodsType models.GoodsType
		err = json.Unmarshal([]byte(result), &goodsType)
		if err != nil {
			logx.GetLogger("ShopManage").Errorf("Handler|GetDeleteTypeList|JsonUnmarshalError|%v", err)
			ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "查询失败", nil))
			ctx.Abort()
			return
		}
		deleteType = append(deleteType, goodsType)
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取商品列表成功", deleteType))
}

func GetTypeInfo(ctx *gin.Context) {
	typeId := ctx.Param("type_id")
	var typrInfo models.GoodsType
	err := database.MyDB.Where("type_id = ?", typeId).First(&typrInfo).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|GetTypeInfo|MySqlError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "该分类不存在", nil))
		ctx.Abort()
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取商品信息成功", typrInfo))
}

// SearchType 搜索商品分类
// @Summary 搜索商品分类
// @Description 搜索商品分类
// @Tags type
// @Accept json
// @Produce json
// @Param type_name query string true "type_name"
// @Router /type/search [get]
func SearchType(ctx *gin.Context) {
	typeName := ctx.Query("type_name")
	if len(typeName) < 0 {
		logx.GetLogger("ShopManage").Errorf("Handler|SearchType|ParamError|%v", "分类名不可以为空")
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "分类名不可以为空", nil))
		ctx.Abort()
	}

	var goodsType []models.GoodsType
	err := database.MyDB.Where("type_name LIKE ?", "%"+typeName+"%").Find(&goodsType).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Handler|SearchType|MySqlError|%v", err)
		panic(err)
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取商品列表成功", goodsType))
}
