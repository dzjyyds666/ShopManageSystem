package models

type GoodsInfo struct {
	GoodsId  string `gorm:"column:goods_id;primaryKey;size:36;not null"`
	Name     string `gorm:"column:name;size:36;not null;unique"`
	Photo    string `gorm:"column:photo;size:256;"`
	Price    int    `gorm:"column:price;size:36;not null"`        // 原价
	Discount int    `gorm:"column:discount;not null;defalut:1.0"` // 折扣
	Stock    int    `gorm:"column:stock;not null"`                // 库存
	TypeId   string `gorm:"column:type_id;size:36;not null" json:"type_id"`
	Status   int    `gorm:"column:status;size:1;not null"`
	// 描述
	Description string `gorm:"column:description;size:300;not null"`
}

func (g GoodsInfo) TableName() string {
	return "goods_info"
}

var GoodsStatus = struct {
	Normal  int
	Deleted int
	Offline int
}{
	Normal:  0,
	Deleted: 1,
	Offline: 2,
}

//type goodsInterface interface {
//	GetGoodPrice(prices ...int) int                           // 获取商品价格
//}
//
//func (g Goods) GetGoodPrice(prices ...int) int {
//	var price int
//	if len(prices) > 0 {
//		for _, item := range prices {
//			price = price + item
//		}
//	}
//	return price
//}
