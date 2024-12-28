package settlement

// todo 策略模式

// 结算方式，会员0.9折
type Strategy interface {
	CalculateTotal(items ...float64) float64
}

// 普通用户支付
type NormalStrategy struct{}

func NewNormalStrategy() NormalStrategy {
	return NormalStrategy{}
}

func (ns NormalStrategy) CalculateTotal(items ...float64) float64 {
	total := 0.0
	for _, item := range items {
		total += item
	}
	return total
}

// 会员支付
type MemberStrategy struct{}

func NewVipStrategy() MemberStrategy {
	return MemberStrategy{}
}

func (ms MemberStrategy) CalculateTotal(items ...float64) float64 {
	total := 0.0
	for _, item := range items {
		total += item
	}
	return total * 0.9 // 会员享受9折
}

// 结算
type Context struct {
	PubStrategy Strategy
}

func (c Context) SetStrategy(strategy Strategy) {
	c.PubStrategy = strategy
}

func (c Context) CalculateTotal(items ...float64) float64 {
	return c.PubStrategy.CalculateTotal(items...)
}
