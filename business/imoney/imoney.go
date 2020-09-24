package imoney

import (
	m_imoney "github.com/gingerxman/ginger-finance/models/imoney"
)

type Imoney struct {
	Code string
	DisplayName string
	ExchangeRate float64
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

func NewImoneyFromModel(dbModel *m_imoney.IMoney) *Imoney{
	instance := new(Imoney)

	instance.Code = dbModel.Code
	instance.DisplayName = dbModel.DisplayName
	instance.ExchangeRate = dbModel.ExchangeRate
	instance.IsPayable = dbModel.IsPayable
	instance.IsDebtable = dbModel.IsDebtable

	return instance
}