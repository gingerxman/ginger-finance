package imoney

import (
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	b_user "github.com/gingerxman/ginger-finance/business/user"
)

type Balance struct {
	eel.RestResource
}

func (this *Balance) Resource() string{
	return "imoney.balance"
}

func (this *Balance) GetParameters() map[string][]string{
	return map[string][]string{
		"GET": []string{"imoney_code", "?with_options:json"},
	}
}

// Get 虚拟资产余额
func (this *Balance) Get(ctx *eel.Context){
	bCtx := ctx.GetBusinessContext()
	req := ctx.Request
	
	imoneyCode := req.GetString("imoney_code")
	user := b_user.GetUserFromContext(bCtx)

	account := b_account.NewAccountRepository(bCtx).GetByUser(user, imoneyCode)
	ctx.Response.JSON(account.Balance)
}