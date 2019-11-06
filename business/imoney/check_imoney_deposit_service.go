package imoney

import (
	"context"
	"encoding/json"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/business/imoney/params"
)


type CheckImoneyDepositService struct{
	eel.ServiceBase
}

func (this *CheckImoneyDepositService) getDepositOrder(bid string) map[string]interface{}{
	resp, err := eel.NewResource(this.Ctx).Get("peanut", "order.deposit_order", eel.Map{
		"bid": bid,
	})

	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("deposit:check_failed", "获取充值订单详情失败"))
	}
	return resp.Data().MustMap()
}

// CheckDepositOrder 校验充值订单
func (this *CheckImoneyDepositService) Check(params params.DepositParams){
	orderData := this.getDepositOrder(params.Bid)
	if orderData["status"].(string) == "finished"{
		panic(eel.NewBusinessError("deposit:check_failed", "不合法的充值订单"))
	}
	if orderData["imoney_code"].(string) != params.ImoneyCode{
		panic(eel.NewBusinessError("deposit:check_failed", "不合法的充值资产"))
	}
	depositAmount, _ := orderData["amount"].(json.Number).Int64()

	if int(depositAmount) != params.Amount {
		panic(eel.NewBusinessError("deposit:check_failed", "不合法的充值金额"))
	}

	depositUserId , _ := orderData["user_id"].(json.Number).Int64()

	if int(depositUserId) != params.UserId{
		panic(eel.NewBusinessError("deposit:check_failed", "不合法的充值用户"))
	}
}

func NewCheckImoneyDepositService(ctx context.Context) *CheckImoneyDepositService{
	instance := new(CheckImoneyDepositService)
	instance.Ctx = ctx
	return instance
}