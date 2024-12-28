package database

import (
	"ShopManageSystem/config"
	"ShopManageSystem/utils/log/logx"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

// todo 单例模式 全局redis客户端
var RDB []*redis.Client

const (
	Redis_Token_Key             = "token:%s"
	Redis_Captcha_Key           = "captcha:%s"
	Redis_Verification_Code_Key = "verification_code:%s"
	Redis_Register_Verify_Key   = "register_verify:%s"

	Redis_GoodS_Stock_Key    = "goods_stock:%s"
	Redis_User_Order_Key     = "user_order:%s"
	Redis_Discount_Goods_Key = "discount_goods_ids"
)

// 初始化redis客户端
func InitRedis() {
	rdb := config.GlobalConfig.Redis
	for _, db := range rdb.DB {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", rdb.Host, rdb.Port),
			Password: rdb.Password,
			DB:       db,
		})
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			logx.GetLogger("SH").Errorf("Database|RedisConnect|FAIL|%v|%v", db, err)
			os.Exit(1)
		} else {
			RDB = append(RDB, client)
		}
	}
	logx.GetLogger("SH").Info("Database|RedisConnect|SUCC")
}
