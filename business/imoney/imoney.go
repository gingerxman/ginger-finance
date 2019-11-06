package imoney

import (
	"context"
	
	"github.com/gingerxman/eel"
	m_imoney "github.com/gingerxman/ginger-finance/models/imoney"
)

type Imoney struct {
	eel.EntityBase

	Id int
	Code string
	DisplayName string
	ExchangeRate float64
	EnableFraction bool
	IsPayable bool
	IsDebtable bool
}

// ExchangeCash 兑换成现金
// 比如：mpcoin的exchange_rate是0.5，则，exchange_cash(10) => 5
// 结果四舍五入
func (this *Imoney) ExchangeCash(amount int) int {
	if amount <= 0 {
		return 0
	}
	
	//TODO: 兑换成现金，改进浮点计算
	return int(this.ExchangeRate * float64(amount))
}

func NewImoneyFromModel(ctx context.Context, dbModel *m_imoney.IMoney) *Imoney{
	instance := new(Imoney)
	instance.Ctx = ctx
	instance.Model = dbModel

	instance.Id = dbModel.Id
	instance.Code = dbModel.Code
	instance.DisplayName = dbModel.DisplayName
	instance.ExchangeRate = dbModel.ExchangeRate
	instance.EnableFraction = dbModel.EnableFraction
	instance.IsPayable = dbModel.IsPayable
	instance.IsDebtable = dbModel.IsDebtable

	return instance
}

func NewImoneyFromMap(ctx context.Context, data map[string]interface{}) *Imoney{
	instance := new(Imoney)
	instance.Ctx = ctx

	code := data["code"].(string)
	displayName := code
	if val, ok := data["display_name"]; ok && val !=nil{
		displayName = val.(string)
	}

	enbaleFraction := true
	if val, ok := data["enbale_fraction"]; ok{
		enbaleFraction = val.(bool)
	}

	isPayable := true
	if val, ok := data["is_payable"]; ok{
		isPayable = val.(bool)
	}

	isDebtable := false
	if val, ok := data["is_debtable"]; ok{
		isDebtable = val.(bool)
	}

	instance.Code = code
	instance.DisplayName = displayName
	instance.ExchangeRate = data["exchange_rate"].(float64)
	instance.EnableFraction = enbaleFraction
	instance.IsPayable = isPayable
	instance.IsDebtable = isDebtable

	return instance
}

// CreateImoney
func CreateImoney(ctx context.Context, data map[string]interface{}) *Imoney{
	imoney := NewImoneyFromMap(ctx, data)
	
	o := eel.GetOrmFromContext(ctx)
	
	if (o.Model(&m_imoney.IMoney{}).Where("code", imoney.Code).Exist()) {
		var model m_imoney.IMoney
		db := o.Model(&m_imoney.IMoney{}).Where("code", imoney.Code).Take(&model)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
		}
		
		return NewImoneyFromModel(ctx, &model)
	} else {
		db := eel.GetOrmFromContext(ctx).Create(&m_imoney.IMoney{
			Code: imoney.Code,
			DisplayName: imoney.DisplayName,
			ExchangeRate: imoney.ExchangeRate,
			EnableFraction: imoney.EnableFraction,
			IsPayable: imoney.IsPayable,
			IsDebtable: imoney.IsDebtable,
		})
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("imoney:create_failed", "创建虚拟资产失败"))
		}
		return imoney
	}
}