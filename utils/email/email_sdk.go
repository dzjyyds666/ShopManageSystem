package email

import (
	"ShopManageSystem/config"
	"ShopManageSystem/utils/log/logx"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(config.GlobalConfig.Email.From, config.GlobalConfig.Email.Alias))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		config.GlobalConfig.Email.Host,
		config.GlobalConfig.Email.Port,
		config.GlobalConfig.Email.User,
		config.GlobalConfig.Email.Password)

	if err := d.DialAndSend(m); err != nil {
		logx.GetLogger("ShopManage").Errorf("发送邮件失败: %v", err)
		return err
	}
	logx.GetLogger("ShopManage").Info("发送邮件成功")
	return nil
}
