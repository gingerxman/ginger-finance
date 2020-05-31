package transfer

import (
	"context"
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	m_account "github.com/gingerxman/ginger-finance/models/account"
	"strings"
	"time"
)

type Transfer struct {
	eel.EntityBase

	Id int
	Bid string
	ThirdBid string
	SourceAccountId int
	DestAccountId int
	SourceAmount int
	DestAmount int
	Action string
	ExtraData string
	Digest string
	IsDeleted bool
	CreatedAt  time.Time

	SourceAccount *b_account.Account
	DestAccount *b_account.Account
}

func (this *Transfer) GetBid() string{
	return this.Bid
}

func (this *Transfer) GetBusinessType() string{
	return "transfer"
}

func (this *Transfer) GetActionType() string{
	if strings.HasPrefix(this.Action, "settlement: bid_"){
		return "settlement"
	}else if strings.HasPrefix(this.Action, "re_settlement: bid_"){
		return "re_settlement"
	}else{
		return this.Action
	}
}

func (this *Transfer) GetDisplayAction() string{
	switch this.Action {
	case "purchase":
		return "购物"
	case "member_purchase":
		return "会员支付"
	case "reverse_transfer":
		return "退款"
	case "imoney.transfer:offline":
		return "线下支付"
	case "imoney.deposit":
		return "会员充值"
	case "imoney.purchase:member_grade":
		return "购买会员"
	case "imoney.withdraw":
		return "提现"
	case "imoney.withdraw:reject":
		return "提现驳回"
	default:
		if strings.HasPrefix(this.Action, "settlement: bid_"){
			return "订单清算"
		}else if strings.HasPrefix(this.Action, "re_settlement: bid_"){
			return "订单重新清算"
		}else if strings.HasPrefix(this.Action, "direct: bid_"){
			return "转账"
		}else{
			return "支付"
		}
	}
}

func (this *Transfer) GetTransferType(account *b_account.Account) string{
	if this.SourceAccountId == account.Id{
		return "expense"
	}else if this.DestAccountId == account.Id{
		return "income"
	}else{
		return "unknown"
	}
}

func NewTransferFromModel(ctx context.Context, dbModel *m_account.Transfer) *Transfer{
	instance := new(Transfer)
	instance.Ctx = ctx
	instance.Model = dbModel

	instance.Id = dbModel.Id
	instance.Bid = dbModel.Bid
	instance.ThirdBid = dbModel.ThirdBid
	instance.SourceAccountId = dbModel.SourceAccountId
	instance.DestAccountId = dbModel.DestAccountId
	instance.SourceAmount = dbModel.SourceAmount
	instance.DestAmount = dbModel.DestAmount
	instance.Action = dbModel.Action
	instance.ExtraData = dbModel.ExtraData
	instance.Digest = dbModel.Digest
	instance.IsDeleted = dbModel.IsDeleted
	instance.CreatedAt = dbModel.CreatedAt

	return instance
}

func NewTransferFromMap(data map[string]interface{}) *Transfer{
	instance := new(Transfer)

	sourceAccount := data["source_account"].(*b_account.Account)
	destAccount := data["dest_account"].(*b_account.Account)

	instance.Bid = data["bid"].(string)
	instance.ThirdBid = data["third_bid"].(string)
	instance.SourceAccountId = sourceAccount.Id
	instance.DestAccountId = destAccount.Id
	instance.SourceAmount = data["source_amount"].(int)
	instance.DestAmount = data["dest_amount"].(int)
	instance.Action = data["action"].(string)
	instance.ExtraData = data["extra_data"].(string)
	instance.Digest = data["digest"].(string)
	instance.CreatedAt = time.Now()

	instance.SourceAccount = sourceAccount
	instance.DestAccount = destAccount

	return instance
}