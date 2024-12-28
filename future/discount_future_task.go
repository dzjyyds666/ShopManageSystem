package future

import (
	"ShopManageSystem/utils/log/logx"
	"time"
)

type Discount struct {
	goodsIds string
	discount float32
}

var discount_task_chan = make(chan Discount, 10)

func StartListenDiscountTask() {
	go UpdateDiscount()
}

func AddDiscountTask(goodsIds string, discount float32) {
	discount_task_chan <- Discount{
		goodsIds: goodsIds,
		discount: discount,
	}
}

func StopListenDiscountTask() {
	close(discount_task_chan)
}

func UpdateDiscount() {
	for {
		select {
		case discount, ok := <-discount_task_chan:
			if !ok {
				logx.GetLogger("ShopManage").Infof("discount_task_chan 已关闭，退出 UpdateDiscount")
				return
			}
			logx.GetLogger("ShopManage").Infof("开始处理任务: %v", discount.goodsIds)
		default:
			// 休眠100毫秒
			time.Sleep(200 * time.Millisecond)
		}
	}
}
