package clearance

import (
	"context"
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
)


type FillClearanceRecordService struct{
	eel.ServiceBase
}

func (this *FillClearanceRecordService) Fill(records []*ClearanceRecord, option eel.FillOption) {
	if option == nil{
		return
	}

	if f, ok := option["with_account"]; ok && f{
		this.fillAccount(records)
	}
}

// fillAccount 填充账户信息
func (this *FillClearanceRecordService) fillAccount(records []*ClearanceRecord){
	if len(records) == 0{
		return
	}
	accountIds := make([]int, 0)
	for _, record := range records{
		accountIds = append(accountIds, record.DestAccountId)
		accountIds = append(accountIds, record.SourceAccountId)
	}
	id2Account := b_account.NewAccountRepository(this.Ctx).GetId2Account(accountIds)

	for _, record := range records{
		record.SourceAccount = id2Account[record.SourceAccountId]
		record.DestAccount = id2Account[record.DestAccountId]
	}
}

func NewFillClearanceRecordService(ctx context.Context) *FillClearanceRecordService{
	instance := new(FillClearanceRecordService)
	instance.Ctx = ctx
	return instance
}