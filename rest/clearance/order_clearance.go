package clearance

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/config"
	"github.com/gingerxman/ginger-finance/business/clearance"
)

type OrderClearance struct {
	eel.RestResource
}

// 对订单号加锁，防止同一笔订单号并发清算
func (this *OrderClearance) GetLockKey(ctx *eel.Context) string {
	req := ctx.Request
	orderBid := req.GetString("bid")
	if orderBid != ""{
		return fmt.Sprintf("order_clearance_%s", orderBid)
	}
	return ""
}

func (this *OrderClearance) Resource() string{
	return "clearance.order_clearance"
}

func (this *OrderClearance) GetParameters() map[string][]string{
	return map[string][]string{
		"PUT": []string{"bid", "?order_status", "?sync:bool", "?_v:int"},
	}
}

// Put 订单清算
func (this *OrderClearance) Put(ctx *eel.Context){
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()

	orderBid := req.GetString("bid")
	if orderBid == ""{
		panic(eel.NewBusinessError("order_clearance:invalid_order_bid", "不合法的订单编号"))
	}
	orderStatus := req.GetString("order_status", "finished")
	sync, _ := req.GetBool("sync", true)
	
	clearanceEngine := clearance.NewClearanceEngine(bCtx)
	if clearanceEngine.IsOrderCleared(orderBid){
		ctx.Response.JSON(eel.Map{
			"status": "success1",
		})
		return
	}
	
	recordIds := clearanceEngine.DoClearance(orderBid, orderStatus)
	settlementMode := config.ServiceConfig.DefaultBool("settlement:RUN_MODE", true)
	if settlementMode {
		sync = true
	}
	
	if sync {
		// do settlement
		spew.Dump(recordIds)
		clearance.NewSettleService(bCtx).DoSettle(recordIds)
	}
	
	ctx.Response.JSON(eel.Map{
		"status": "success2",
	})
}