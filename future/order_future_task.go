package future

import (
	"ShopManageSystem/utils/log/logx"
	"time"
)

type OrderInfo struct {
	OrderId string
	number  int
}

var order_task_chan = make(chan OrderInfo, 10)

func StartListenOrderTask() {
	go UpdateGoodsStock()
}

func AddOrderTask(orderId string, number int) {
	orderinfo := OrderInfo{
		OrderId: orderId,
		number:  number,
	}
	order_task_chan <- orderinfo
}

func StopListenOrderTask() {
	close(order_task_chan)
}

func UpdateGoodsStock() {
	for {
		select {
		case task, ok := <-order_task_chan:
			if !ok {
				logx.GetLogger("ShopManage").Infof("order_task_chan 已关闭，退出 UpdateGoodsStock")
				return
			}
			logx.GetLogger("ShopManage").Infof("开始处理任务: %v", task.OrderId)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
