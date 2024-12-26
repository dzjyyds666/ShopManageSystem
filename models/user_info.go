package models

type UserInfo struct {
	UserId   string `gorm:"column:user_id;primaryKey;size:36;not null"`
	UserName string `gorm:"column:user_name;size:36;not null;unique"`
	Email    string `gorm:"column:email;size:36;not null,unique"`
	Password string `gorm:"column:password;size:36;not null"`
	Avatar   string `gorm:"column:avatar;size:36;not null"`
}
