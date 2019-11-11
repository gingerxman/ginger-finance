package account

import (
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	b_user "github.com/gingerxman/ginger-finance/business/user"
)

type Account struct {
	eel.RestResource
}

func (this *Account) Resource() string{
	return "account.account"
}

func (this *Account) SkipAuthCheck() bool {
	return true
}

func (this *Account) GetParameters() map[string][]string{
	return map[string][]string{
		"GET": []string{},
		"PUT": []string{"user_id:int", "imoney_code", "is_debtable:bool"},
	}
}

func (this *Account) Get(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	corp := b_user.NewCorpFromOnlyId(bCtx, 3358)
	
	ctx.Response.JSON(eel.Map{
		"user_id": corp.GetRelatedUser().Id,
	})
}

func (this *Account) Put(ctx *eel.Context){
	bCtx := ctx.GetBusinessContext()
	req := ctx.Request
	userId, _ := req.GetInt("user_id", 0)
	imoneyCode := req.GetString("imoney_code", "")
	isDebtable, _ := req.GetBool("is_debtable", true)
	
	if userId == 0{
		panic(eel.NewBusinessError("invalid_user_id", "不合法的参数:user_id(0)"))
	}
	user := b_user.NewUserFromOnlyId(bCtx, userId)
	
	b_account.NewAccountFactory(bCtx).CreateForUser(user, b_account.CreateAccountParams{
		ImoneyCode: imoneyCode,
		IsDebtable: isDebtable,
		WithDebtable: true,
	})
	
	ctx.Response.JSON(eel.Map{})
}