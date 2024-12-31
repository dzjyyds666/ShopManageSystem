package models

type GoodsType struct {
	TypeId      string `gorm:"column:type_id;primaryKey;size:36;not null"`
	GoodsNumber string `gorm:"column:goods_number;size:36;not null;default:0"`
	TypeName    string `gorm:"column:type_name;size:36;not null;unique"`
}

func (g GoodsType) TableName() string {
	return "goods_type"
}
