package imoney

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/business/account/frozen"
	b_user "github.com/gingerxman/ginger-finance/business/user"
	b_account "github.com/gingerxman/ginger-finance/business/account"
)

type FrozenRecord struct {
	eel.RestResource
}

func (this *FrozenRecord) Resource() string{
	return "imoney.frozen_record"
}

// GetLockKey 加锁
func (this *FrozenRecord) GetLockKey(ctx *eel.Context) string {
	bCtx := ctx.GetBusinessContext()
	user := b_user.GetUserFromContext(bCtx)
	return fmt.Sprintf("imoney_frozen_record_%d", user.GetId())
}

func (this *FrozenRecord) GetParameters() map[string][]string{
	return map[string][]string{
		"GET": []string{"imoney_code"},
		"PUT": []string{"imoney_code", "amount:int", "type", "?remark"},
		"DELETE": []string{"id:int"},
	}
}

func (this *FrozenRecord) Get(ctx *eel.Context){
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	imoneyCode := req.GetString("imoney_code")

	user := b_user.GetUserFromContext(bCtx)

	account := b_account.NewAccountRepository(bCtx).GetByUser(user, imoneyCode)
	amount := account.FrozenAmount
	
	ctx.Response.JSON(eel.Map{
		"frozen_amount": amount,
	})
}

// Put 冻结资产
func (this *FrozenRecord) Put(ctx *eel.Context){
	
	req := ctx.Request
	
	imoneyCode := req.GetString("imoney_code")
	remark := req.GetString("remark", "")
	frozenType := req.GetString("type")
	amount, _ := req.GetInt("amount")
	
	bCtx := ctx.GetBusinessContext()
	user := b_user.GetUserFromContext(bCtx)
	account := b_account.NewAccountRepository(bCtx).GetByUser(user, imoneyCode)

	record := frozen.NewFrozenAccountService(bCtx).FrozenForAccount(account, amount, frozenType, remark)
	ctx.Response.JSON(eel.Map{
		"frozen_record_id": record.Id,
	})
}

func (this *FrozenRecord) Delete(ctx *eel.Context){
	req := ctx.Request
	
	recordId, _ := req.GetInt("id")
	
	bCtx := ctx.GetBusinessContext()
	frozenRecord := frozen.NewFrozenRecordRepository(bCtx).GetById(recordId)
	if frozenRecord == nil {
		ctx.Response.Error("frozen_record:invalid_frozen_record", fmt.Sprintf("%d", recordId))
		return
	}
	
	account := b_account.NewAccountRepository(bCtx).GetById(frozenRecord.AccountId)
	if account == nil {
		ctx.Response.Error("frozen_record:invalid_record_account", fmt.Sprintf("%d-%d", frozenRecord, frozenRecord.AccountId))
		return
	}
	
	frozen.NewFrozenAccountService(bCtx).UnfrozenForAccount(account, frozenRecord)
	ctx.Response.JSON(eel.Map{
	})
}