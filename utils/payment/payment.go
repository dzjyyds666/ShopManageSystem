package payment

import "ShopManageSystem/utils/log/logx"

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
	Pay(amount float64) (bool, error)
}

type WechatPayment struct{}

func (w *WechatPayment) Pay(amount float64) (bool, error) {
	logx.GetLogger("ShopManage").Infof("Payment|WechatPayment|Pay|%d", amount)
	return true, nil
}

type AliPayment struct{}

func (a *AliPayment) Pay(amount float64) (bool, error) {
	logx.GetLogger("ShopManage").Infof("Payment|AliPayment|Pay|%d", amount)
	return true, nil
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
