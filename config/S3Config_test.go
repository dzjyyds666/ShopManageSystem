package config

import (
	"ShopManageSystem/utils/log/logx"
	"testing"
)

func TestLoggingConfig(t *testing.T) {
	config := LoadS3Config("s3_config.toml")
	logx.GetLogger("ShopManage").Infof("Config|LoadS3Config|%v", config)
}
