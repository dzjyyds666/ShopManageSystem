package handlers

import (
	"ShopManageSystem/database"
	"ShopManageSystem/future"
	"ShopManageSystem/models"
	"ShopManageSystem/utils/log/logx"
	"ShopManageSystem/utils/payment"
	"ShopManageSystem/utils/response"
	"ShopManageSystem/utils/settlement"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

const (
	Order_Expire_Time = 60 * 10 // 订单过期时间
)

// @Summary 分页获取商品列表
// @Description 分页获取商品列表
// @Tags goods
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param limit query int true "每页数量"
// @Router /goods/list [get]
func GetGoodsListByPage(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	offset := (pageInt - 1) * limitInt

	var goodsList []models.GoodsInfo
	err := database.MyDB.Where("status = ?", models.GoodsStatus.Normal).
		Offset(offset).
		Limit(limitInt).
		Select("goods_id", "name", "original_price", "real_price", "stock", "type_id", "photo").
		Find(&goodsList).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("GetGoodsListByPage|MySqlError|%v", err)
		panic(err)
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取商品列表成功", goodsList))
}

// @Summary 获取单个商品信息
// @Description 获取单个商品信息
// @Tags goods
// @Accept json
// @Produce json
// @Param goods_id path string true "商品id"
// @Router /goods/info/{goods_id} [get]
func GetGoodsInfo(ctx *gin.Context) {
	goodId := ctx.Param("goods_id")

	if len(goodId) < 0 {
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "货物id不可以为空", nil))
		ctx.Abort()
	}

	var goodsInfo models.GoodsInfo
	err := database.MyDB.Where("goods_id = ?", goodId).First(&goodsInfo).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("GetGoodsInfo|MySqlError|%v", err)
		panic(err)
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "获取商品信息成功", goodsInfo))
}

// @Summary 上架商品
// @Description 上架商品
// @Tags goods
// @Accept json
// @Produce json
// @Param goods_info body models.GoodsInfo true "商品信息"
// @Router /goods/onShelves [post]
func GoodsOnShelves(ctx *gin.Context) {
	var goodinfo models.GoodsInfo
	err := ctx.ShouldBindJSON(&goodinfo)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("GoodsOnShelves|ParamsError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "参数错误", nil))
		ctx.Abort()
	}

	goodIdUUID, _ := uuid.NewUUID()
	goodId := strings.ReplaceAll(goodIdUUID.String(), "-", "")
	goodinfo.GoodsId = goodId
	// 把商品信息存入mysql中
	err = database.MyDB.Create(&goodinfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) { // 判断是否重复
			logx.GetLogger("ShopManage").Errorf("GoodsOnShelves|DuplicateKeyError|%v", err)
			ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "该商品已存在", nil))
			ctx.Abort()
		}
		logx.GetLogger("ShopManage").Errorf("GoodsOnShelves|MySqlError|%v", err)
		panic(err)
	}
	// 把商品的库存信息存入redis中
	err = database.RDB[0].Set(ctx, fmt.Sprintf(database.Redis_GoodS_Stock_Key, goodId), goodinfo.Stock, 0).Err()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("GoodsOnShelves|RedisError|%v", err)
		panic(err)
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "上架成功", nil))
}

type buyGoods struct {
	GoodsId string    `json:"goods_id"`
	Num     int       `json:"num"`
	Prices  []float32 `json:"prices"`
}

type orderInfo struct {
	Time   string    `json:"time"`
	Price  []float64 `json:"Price"`
	Number int       `json:"number"`
	GoodId string    `json:"good_id"`
	UserId string    `json:"user_id"`
	Role   string    `json:"role"`
}

// InitBuyGoods 初始化订单
// @Summary 下订单
// @Description 下订单
// @Tags goods
// @Accept json
// @Produce json
// @Param buy_goods body handlers.buyGoods true "购买商品信息"
// @Router /goods/initPayment [post]
func InitBuyGoods(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	var buygoods buyGoods
	err := ctx.ShouldBindJSON(&buygoods)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("BuyGoods|ParamsError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "参数错误", nil))
		ctx.Abort()
	}

	get := database.RDB[0].Get(ctx, fmt.Sprintf(database.Redis_GoodS_Stock_Key, buygoods.GoodsId))
	if get.Err() != nil {
		logx.GetLogger("ShopManage").Errorf("BuyGoods|GetStockError|%v", err)
		panic(get.Err())
	}

	stock, _ := strconv.Atoi(get.Val())
	if stock < buygoods.Num {
		logx.GetLogger("ShopManage").Errorf("BuyGoods|StockNotEnough|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "库存不足，购买失败", nil))
		ctx.Abort()
	}

	//减去redis的订单信息
	stockkey := fmt.Sprintf(database.Redis_GoodS_Stock_Key, buygoods.GoodsId)

	// Lua 脚本
	luaScript := `
	local stockKey = KEYS[1]
	local decrement = tonumber(ARGV[1])

	-- 获取当前库存
	local stock = tonumber(redis.call("GET", stockKey))

	-- 检查库存是否足够
	if stock and stock >= decrement then
		-- 减少库存
		redis.call("DECRBY", stockKey, decrement)
		return true
	else
		return false
	end
	`

	result, err := database.RDB[0].Eval(ctx, luaScript, []string{stockkey}, buygoods.Num).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("BuyGoods|RedisError|luaScript|%v", err)
		panic(err)
	}

	if result == false {
		logx.GetLogger("ShopManage").Errorf("BuyGoods|StockNotEnough|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "库存不足，购买失败", nil))
		ctx.Abort()
	}

	// 查询用户身份信息
	var userInfo models.UserInfo
	err = database.MyDB.Where("user_id = ?", userId).First(&userInfo).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("BuyGoods|MySqlError|%v", err)
		panic(err)
	}

	orderId := GenerateVerificationCode(10)

	orderinfo := orderInfo{
		Time: strconv.FormatInt(time.Now().Unix(), 10),

		Number: buygoods.Num,
		GoodId: buygoods.GoodsId,
		UserId: userId.(string),
		Role:   userInfo.Role,
	}

	rawData, _ := json.Marshal(&orderinfo)

	err = database.RDB[0].Set(ctx, fmt.Sprintf(database.Redis_User_Order_Key, orderId), rawData, 0).Err()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("BuyGoods|RedisError|SetOrderError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "初始化订单失败", nil))
		// 还原库存信息
		err = database.RDB[0].IncrBy(ctx, fmt.Sprintf(database.Redis_GoodS_Stock_Key, buygoods.GoodsId), int64(buygoods.Num)).Err()
		if err != nil {
			logx.GetLogger("ShopManage").Errorf("BuyGoods|RedisError|IncrByError|%v", err)
			panic(err)
		}
		ctx.Abort()
	}

	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "初始化订单完成", orderinfo))
}

// PayForOrder 付款
// @Summary 付款
// @Description 付款
// @Tags goods
// @Accept json
// @Produce json
// @Param order_id path string true "订单id"
// @Router /goods/payForOrder/{order_id} [post]
func PayForOrder(ctx *gin.Context) {
	orderId := ctx.Param("order_id")
	if orderId == "" {
		logx.GetLogger("ShopManage").Errorf("PayForOrder|ParamsError|Order Can not be null")
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "订单id不可以为空", nil))
		ctx.Abort()
	}

	payType := ctx.Query("pay_type")

	// 从redis中获取到订单信息
	orderinfo, err := database.RDB[0].Get(ctx, fmt.Sprintf(database.Redis_User_Order_Key, orderId)).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("PayForOrder|RedisError|GetOrderError|%v", err)
		panic(err)
	}

	var order orderInfo
	err = json.Unmarshal([]byte(orderinfo), &order)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("PayForOrder|JsonUnmarshalError|%v", err)
		panic(err)
	}

	// 先判断订单是否过期
	createTime, err := strconv.ParseInt(order.Time, 10, 64)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("PayForOrder|ParseIntError|%v", err)
		panic(err)
	}

	if time.Now().Unix()-createTime > Order_Expire_Time {
		logx.GetLogger("ShopManage").Errorf("PayForOrder|OrderExpired|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.ParamError, "订单已过期", nil))
		ctx.Abort()
	}

	var strategy settlement.Strategy
	if order.Role == models.Role.NormalUser {
		strategy = settlement.NewNormalStrategy()
	} else if order.Role == models.Role.VipUser {
		strategy = settlement.NewVipStrategy()
	}

	// todo 策略模式计算总价格
	var context settlement.Context
	context.SetStrategy(strategy)
	totalPrice := context.PubStrategy.CalculateTotal(order.Price...)

	// todo 简单工厂方法模式
	var payclient payment.Payment
	paymentFactoty := payment.PaymentFactory{}
	payclient = paymentFactoty.CreatePayment(payment.PayType(payType))
	pay, err := payclient.Pay(totalPrice)
	if err != nil {
		return
	}

	if pay == false {
		logx.GetLogger("ShopManage").Errorf("PayForOrder|PaymentError|%v", err)
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.PayFailed, "支付失败", nil))
		ctx.Abort()
	} else {
		logx.GetLogger("ShopManage").Infof("PayForOrder|PaymentSuccess")
		ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "支付成功", nil))
	}
}

// @Summary 完成订单
// @Description 完成订单
// @Tags goods
// @Accept json
// @Produce json
// @Param order_id path string true "订单id"
// @Router /goods/completeOrder/{order_id} [get]
func CompleteOrder(ctx *gin.Context) {
	orderId := ctx.Param("order_id")
	// 删除redis中的订单信息的过期时间，代表完成订单
	result, err := database.RDB[0].Get(ctx, fmt.Sprintf(database.Redis_User_Order_Key, orderId)).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("CompleteOrder|RedisError|%v", err)
		panic(err)
	}

	var order orderInfo
	json.Unmarshal([]byte(result), &order)

	order.Time = ""

	rawData, _ := json.Marshal(&order)
	err = database.RDB[0].Set(ctx, fmt.Sprintf(database.Redis_User_Order_Key, orderId), rawData, 0).Err()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("CompleteOrder|RedisError|%v", err)
		panic(err)
	}

	// 使用异步任务处理订单
	future.AddOrderTask(orderId, order.Number)
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "订单完成", nil))
}

// @Summary 标记打折商品
// @Description 标记打折商品
// @Tags goods
// @Accept json
// @Produce json
// @Param goods_ids query []string true "打折商品id"
// @Router /goods/markDiscountGoods [get]
func MarkDiscountGoods(ctx *gin.Context) {
	goodsIds := ctx.QueryArray("goods_ids")
	for _, goodsId := range goodsIds {
		err := database.RDB[0].SAdd(ctx, fmt.Sprintf(database.Redis_Discount_Goods_Key, goodsId), 1, 0).Err()
		if err != nil {
			logx.GetLogger("ShopManage").Errorf("MarkDiscountGoods|RedisError|%v", err)
			panic(err)
		}
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "打折商品标记成功", nil))
}

// @Summary 取消打折商品
// @Description 取消打折商品
// @Tags goods
// @Accept json
// @Produce json
// @Param goods_ids query []string true "打折商品id"
// @Router /goods/cancelDiscountGoods [get]
func CancelDiscountGoods(ctx *gin.Context) {
	goodsIds := ctx.QueryArray("goods_ids")
	for _, goodsId := range goodsIds {
		err := database.RDB[0].SRem(ctx, fmt.Sprintf(database.Redis_Discount_Goods_Key, goodsId), 1, 0).Err()
		if err != nil {
			logx.GetLogger("ShopManage").Errorf("CancelDiscountGoods|RedisError|%v", err)
			panic(err)
		}
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "打折商品取消成功", nil))
}

// @Summary 设置打折
// @Description 设置打折
// @Tags goods
// @Accept json
// @Produce json
// @Param discount path string true "打折"
// @Router /goods/discount/{discount} [get]
func Discount(ctx *gin.Context) {
	discount := ctx.Param("discount")
	goodsIds, err := database.RDB[0].SMembers(ctx, database.Redis_Discount_Goods_Key).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Discount|RedisError|%v", err)
		panic(err)
	}

	err = database.MyDB.Where("id in (?)", goodsIds).UpdateColumn("discount", discount).Error
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Discount|DBError|%v", err)
		panic(err)
	}
	ctx.JSON(200, response.NewResult(response.EnmuHttptatus.RequestSuccess, "打折成功", nil))
}
