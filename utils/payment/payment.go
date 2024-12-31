package payment

import (
	"ShopManageSystem/database"
	"ShopManageSystem/utils/log/logx"
	"context"
	"strconv"
)

// 简单工厂创建支付方式

type PayType string

func (p PayType) Is(payType PayType) bool {
	return p == payType
}

var PaymentType = struct {
	WechatPayment PayType
	AliPayment    PayType
}{
	WechatPayment: "wechat",
	AliPayment:    "alipay",
}

// todo 简单工厂方法模式实现
type Payment interface {
	Pay(amount float64) (string, error)
	GetPayToTal() (string, error)
}

type WechatPayment struct{}

func (w *WechatPayment) Pay(amount float64) (string, error) {
	logx.GetLogger("ShopManage").Infof("Payment|WechatPayment|Pay|%d", amount)
	// 存入微信支付的redis中
	result, err := database.RDB[0].Get(context.TODO(), database.Redis_Wechat_Pay_ToTal).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|WechatPayment|Rediserror|%s", err)
		return "微信支付失败", err
	}
	total, err := strconv.ParseFloat(result, 64)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|WechatPayment|AtoiError|%s", err)
		return "微信支付失败", err
	}

	total = total + amount

	err = database.RDB[0].Set(context.TODO(), database.Redis_Wechat_Pay_ToTal, total, 0).Err()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|WechatPayment|RedisError|SetError|%s", err)
		return "微信支付失败", err
	}

	return "微信支付成功", nil
}

func (w *WechatPayment) GetPayToTal() (string, error) {
	result, err := database.RDB[0].Get(context.TODO(), database.Redis_Wechat_Pay_ToTal).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|WechatPayment|Rediserror|%s", err)
		return "获取微信支付总额失败", err
	}
	return result, nil
}

type AliPayment struct{}

func (a *AliPayment) Pay(amount float64) (string, error) {
	logx.GetLogger("ShopManage").Infof("Payment|AliPayment|Pay|%d", amount)
	// 存入支付宝支付的redis中
	result, err := database.RDB[0].Get(context.TODO(), database.Redis_Ali_Pay_ToTal).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|AliPayment|Rediserror|%s", err)
		return "支付宝支付失败", err
	}
	total, err := strconv.ParseFloat(result, 64)
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|AliPayment|AtoiError|%s", err)
		return "支付宝支付失败", err
	}

	total = total + amount

	err = database.RDB[0].Set(context.TODO(), database.Redis_Ali_Pay_ToTal, total, 0).Err()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|AliPayment|RedisError|SetError|%s", err)
		return "支付宝支付失败", err
	}
	return "支付宝支付成功", nil
}

func (a *AliPayment) GetPayToTal() (string, error) {
	result, err := database.RDB[0].Get(context.TODO(), database.Redis_Ali_Pay_ToTal).Result()
	if err != nil {
		logx.GetLogger("ShopManage").Errorf("Payment|AliyunPayment|Rediserror|%s", err)
		return "获取支付宝支付总额失败", err
	}
	return result, nil
}

type PaymentFactory struct{}

func (p *PaymentFactory) CreatePayment(paymentType PayType) Payment {
	switch paymentType {
	case PaymentType.WechatPayment:
		return &WechatPayment{}
	case PaymentType.AliPayment:
		return &AliPayment{}
	default:
		return &WechatPayment{}
	}
}
