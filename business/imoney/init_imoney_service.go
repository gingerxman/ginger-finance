package imoney

import (
	"context"
	
	"github.com/gingerxman/eel"
	m_imoney "github.com/gingerxman/ginger-finance/models/imoney"
)

var Code2Imoney = map[string]map[string]interface{}{
	"rmb":  map[string]interface{}{
		"code": "rmb",
		"exchange_rate": 1.0,
		"enbale_fraction": true,
		"is_debtable": true,
	},
	"cash": map[string]interface{}{
		"code": "cash",
		"exchange_rate": 1.0,
		"enbale_fraction": true,
	},
	"cash_fee": map[string]interface{}{
		"code": "cash_fee", // 手续费
		"exchange_rate": 1.0,
		"enbale_fraction": true,
	},
	"withdraw_cash": map[string]interface{}{
		"code": "withdraw_cash",
		"exchange_rate": 1.0,
		"enbale_fraction": true,
	},
}

type InitImoneyService struct{
	eel.ServiceBase
}

func (this *InitImoneyService) Init() {
	o := eel.GetOrmFromContext(this.Ctx)
	// 首先清空数据
	db := o.Delete(&m_imoney.IMoney{})
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("init_imoney:failed", "清空数据失败"))
	}

	for _, imoney := range Code2Imoney{
		isDebtable := false
		if val, ok := imoney["is_debtable"]; ok{
			isDebtable = val.(bool)
		}
		dbModel := &m_imoney.IMoney{
			Code: imoney["code"].(string),
			DisplayName: imoney["code"].(string),
			ExchangeRate: imoney["exchange_rate"].(float64),
			EnableFraction: imoney["enbale_fraction"].(bool),
			IsDebtable: isDebtable,
			IsPayable: true,
		}
		db = o.Create(dbModel)
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("init_imoney:failed", "保存数据失败"))
		}
	}
}

func NewInitImoneyService(ctx context.Context) *InitImoneyService{
	instance := new(InitImoneyService)
	instance.Ctx = ctx
	return instance
}