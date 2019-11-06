package imoney

import (
	"context"
	"github.com/gingerxman/eel"
)


type ImoneyExchangeService struct{
	eel.ServiceBase
}

// ExchangeByImoneyCode 虚拟资产兑换
func (this *ImoneyExchangeService) ExchangeByImoneyCode(sourceImoneyCode, destImoneyCode string, sourceAmount int) int {
	if sourceImoneyCode == destImoneyCode{
		return sourceAmount
	}
	sourceImoney := NewImoneyRepository(this.Ctx).GetByCode(sourceImoneyCode)
	destImoney := NewImoneyRepository(this.Ctx).GetByCode(destImoneyCode)
	return this.Exchange(sourceImoney, destImoney, sourceAmount)
}

func (this *ImoneyExchangeService) Exchange(sourceImoney, destImoney *Imoney, amount int) int {
	sourceCashAmount := sourceImoney.ExchangeCash(amount)
	
	return int(float64(sourceCashAmount) * destImoney.ExchangeRate)
	// return sourceCashAmount.Mul(decimal.NewFromFloat(destImoney.ExchangeRate)).Round(2)
}

// ExchangeCash 兑换成现金数额
func (this *ImoneyExchangeService) ExchangeCash(sourceImoneyCode string, amount int) int {
	return this.ExchangeByImoneyCode(sourceImoneyCode, "cash", amount)
}

func NewImoneyExchangeService(ctx context.Context) *ImoneyExchangeService{
	instance := new(ImoneyExchangeService)
	instance.Ctx = ctx
	return instance
}