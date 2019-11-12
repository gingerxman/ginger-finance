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
		"GET": []string{"imoney_code", "?view_corp_account:bool", "?with_options:json"},
	}
}

// Get 虚拟资产余额
func (this *Balance) Get(ctx *eel.Context){
	bCtx := ctx.GetBusinessContext()
	req := ctx.Request
	
	viewCorpAccount, _ := req.GetBool("view_corp_account", false)
	imoneyCode := req.GetString("imoney_code")
	
	var user *b_user.User
	if viewCorpAccount {
		corp := b_user.GetCorpFromContext(bCtx)
		user = corp.GetRelatedUser()
	} else {
		user = b_user.GetUserFromContext(bCtx)
	}

	account := b_account.NewAccountRepository(bCtx).GetByUser(user, imoneyCode)
	ctx.Response.JSON(account.Balance)
}