package imoney

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/business/account"
	b_transfer "github.com/gingerxman/ginger-finance/business/transfer"
	b_user "github.com/gingerxman/ginger-finance/business/user"
)

type Transfer struct {
	eel.RestResource
}

func (this *Transfer) Resource() string{
	return "imoney.transfer"
}

// GetLockKey 加锁
func (this *Transfer) GetLockKey(ctx *eel.Context) string {
	bCtx := ctx.GetBusinessContext()
	user := b_user.GetUserFromContext(bCtx)

	return fmt.Sprintf("imoney_transfer_%d", user.GetId())
}

func (this *Transfer) GetParameters() map[string][]string{
	return map[string][]string{
		"PUT": []string{"amount:int", "bid", "source_user_id:int", "dest_user_id:int", "?source_imoney_code", "?dest_imoney_code",
			"imoney_code", "?source_account_code", "?dest_account_code", "?remark",
		},
	}
}

// Put
func (this *Transfer) Put(ctx *eel.Context){
	bCtx := ctx.GetBusinessContext()
	req := ctx.Request
	
	var (
		sourceAccount *account.Account
		destAccount *account.Account
	)
	accountRepository := account.NewAccountRepository(bCtx)
	sourceImoneyCode := req.GetString("imoney_code", "")
	destImoneyCode := req.GetString("imoney_code", "")
	imoneyCode := req.GetString("imoney_code", "")

	//获取source account
	sourceUserId, _ := req.GetInt("source_user_id", 0)
	sourceAccountCode := req.GetString("source_account_code", "")
	if sourceAccountCode != ""{
		sourceAccount = accountRepository.GetByCode(sourceAccountCode)
	}else{
		if sourceImoneyCode == ""{
			sourceImoneyCode = imoneyCode
		}
		sourceAccount = accountRepository.GetByUserId(sourceUserId, sourceImoneyCode)
	}
	
	//dest account
	destUserId, _ := req.GetInt("dest_user_id", 0)
	destAccountCode := req.GetString("dest_account_code", "")
	if destAccountCode != ""{
		destAccount = accountRepository.GetByCode(destAccountCode)
	}else{
		if destImoneyCode == ""{
			destImoneyCode = imoneyCode
		}
		destAccount = accountRepository.GetByUserId(destUserId, destImoneyCode)
	}

	amount, _ := req.GetInt("amount", 0)
	remark := req.GetString("trigger", "直接转账")
	bid := req.GetString("bid")

	transfer := b_transfer.NewTransferService(bCtx).Transfer(b_transfer.TransferParams{
		SourceAccount: sourceAccount,
		DestAccount: destAccount,
		SourceAmount: amount,
		DestAmount: amount,
		Bid: bid,
		Action: fmt.Sprintf("direct: bid_%s", bid),
		ExtraData: map[string]interface{}{
			"remark": remark,
		},
	})
	b_transfer.NewFillTransferService(bCtx).FillId(transfer)

	ctx.Response.JSON(eel.Map{
		"transfer_id": transfer.Id,
		"source_account_balance": sourceAccount.GetValidBalance(),
		"dest_account_balance": destAccount.GetValidBalance(),
	})
}