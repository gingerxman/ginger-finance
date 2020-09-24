package imoney

import (
	"github.com/gingerxman/eel"
	b_imoney "github.com/gingerxman/ginger-finance/business/imoney"
)

type Imoney struct {
	eel.RestResource
}

func (this *Imoney) Resource() string{
	return "imoney.imoney"
}

func (this *Imoney) GetParameters() map[string][]string{
	return map[string][]string{
		"GET": []string{"code"},
		"PUT": []string{"code", "exchange_rate:float", "display_name", "is_payable:bool", "is_debtable:bool"},
	}
}

// Get 获取虚拟资产
func (this *Imoney) Get(ctx *eel.Context){
	bCtx := ctx.GetBusinessContext()
	req := ctx.Request

	imoney := b_imoney.NewImoneyRepository(bCtx).GetByCode(req.GetString("code"))

	ctx.Response.JSON(b_imoney.NewEncodeImoneyService(bCtx).Encode(imoney))
}

func (this *Imoney) Put(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()

	req := ctx.Request
	rate, _ := req.GetFloat("exchange_rate")
	isPayable, _ := req.GetBool("is_payable")
	isDebtable, _ := req.GetBool("is_debtable")

	b_imoney.NewImoneyManager(bCtx).Add(&b_imoney.Imoney{
		Code: req.GetString("code"),
		DisplayName: req.GetString("display_name"),
		ExchangeRate: rate,
		IsPayable: isPayable,
		IsDebtable: isDebtable,
	})
	ctx.Response.JSON(nil)
}