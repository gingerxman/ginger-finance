package transfer

import (
	"context"
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	"github.com/gingerxman/ginger-finance/business/constant"
)


type EncodeTransferService struct{
	eel.ServiceBase
}

func (this *EncodeTransferService) EncodeForAccount(account *b_account.Account, transfer *Transfer) *RTransfer{
	return &RTransfer{
		Id: transfer.Id,
		Bid: transfer.Bid,
		ThirdBid: transfer.ThirdBid,
		SourceAccountId: transfer.SourceAccountId,
		DestAccountId: transfer.DestAccountId,
		SourceUserId: transfer.SourceAccount.UserId,
		DestUserId: transfer.DestAccount.UserId,
		Amount: transfer.SourceAmount,
		SourceAmount: transfer.SourceAmount,
		DestAmount: transfer.DestAmount,
		Action: transfer.GetDisplayAction(),
		OriginAction: transfer.Action,
		TransferType: transfer.GetTransferType(account),
		CreatedAt: transfer.CreatedAt.Format(constant.SHORT_TIME_LAYOUT),
	}
}

func (this *EncodeTransferService) EncodeManyForAccount(account *b_account.Account, transfers []*Transfer) []*RTransfer{
	rTransfers := make([]*RTransfer, 0)
	for _, transfer := range transfers{
		rTransfers = append(rTransfers, this.EncodeForAccount(account, transfer))
	}
	return rTransfers
}

func NewEncodeTransferService(ctx context.Context) *EncodeTransferService{
	instance := new(EncodeTransferService)
	instance.Ctx = ctx
	return instance
}