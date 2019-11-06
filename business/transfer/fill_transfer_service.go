package transfer

import (
	"context"
	mapset "github.com/deckarep/golang-set"
	
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	m_account "github.com/gingerxman/ginger-finance/models/account"
)


type FillTransferService struct{
	eel.ServiceBase
}

func (this *FillTransferService) FillAccount(transfers []*Transfer){

	if len(transfers) == 0{
		return
	}

	sAccountIds := mapset.NewSet()
	for _, transfer := range transfers{
		sAccountIds.Add(transfer.SourceAccountId)
		sAccountIds.Add(transfer.DestAccountId)
	}
	accountIds := make([]int, 0)
	for _, sId := range sAccountIds.ToSlice(){
		accountIds = append(accountIds, sId.(int))
	}
	accounts := b_account.NewAccountRepository(this.Ctx).GetByIds(accountIds)
	id2Account := make(map[int]*b_account.Account)
	for _, account := range accounts{
		id2Account[account.Id] = account
	}

	for _, transfer := range transfers{
		transfer.SourceAccount = id2Account[transfer.SourceAccountId]
		transfer.DestAccount = id2Account[transfer.DestAccountId]
	}
}

// FillId 填充id
func (this *FillTransferService) FillId(transfer *Transfer){
	var transferDbModel m_account.Transfer
	db := eel.GetOrmFromContext(this.Ctx).Model(&m_account.Transfer{}).Where("Bid", transfer.Bid).Take(&transferDbModel)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("fill_transfer:failed", "查询交易失败"))
	}
	transfer.Id = transferDbModel.Id
}

func NewFillTransferService(ctx context.Context) *FillTransferService{
	instance := new(FillTransferService)
	instance.Ctx = ctx
	return instance
}