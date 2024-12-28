package models

type UserInfo struct {
	UserId   string `gorm:"column:user_id;primaryKey;size:36;not null"`
	UserName string `gorm:"column:user_name;size:36;not null;unique"`
	Email    string `gorm:"column:email;size:36;not null,unique"`
	Password string `gorm:"column:password;size:256;not null"`
	Role     string `gorm:"column:role;size:16;not null;default:normal_user"`
	Avatar   string `gorm:"column:avatar;size:36;not null"`
}

func (u UserInfo) TableName() string {
	return "user_info"
}

var Role = struct {
	NormalUser string
	VipUser    string
}{
	NormalUser: "normal_user",
	VipUser:    "vip_user",
}
