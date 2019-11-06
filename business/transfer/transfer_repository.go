package transfer

import (
	"context"
	
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	m_account "github.com/gingerxman/ginger-finance/models/account"
)


type TransferRepository struct{
	eel.ServiceBase
}

func (this *TransferRepository) GetPagedTransfersForAccount(account *b_account.Account, transferType string, filters eel.Map, pageInfo *eel.PageInfo) ([]*Transfer, eel.INextPageInfo){

	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_account.Transfer{})

	switch transferType {
	case "all":
		db = db.Where("source_account_id", account.Id).Or("dest_account_id", account.Id)
	case "income":
		db = db.Where("dest_account_id", account.Id)
	case "expense":
		db = db.Where("source_account_id", account.Id)
	}

	var dbModels []*m_account.Transfer
	db = db.Where(filters).Order("id desc")

	transfers := make([]*Transfer, 0)
	paginateResult, db := eel.Paginate(db, pageInfo, &dbModels)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return transfers, paginateResult
	}

	for _, dbModel := range dbModels {
		transfers = append(transfers, NewTransferFromModel(this.Ctx, dbModel))
	}
	return transfers, paginateResult
}

func (this *TransferRepository) GetByFilters(filters eel.Map) []*Transfer{
	var dbModels []*m_account.Transfer
	db := eel.GetOrmFromContext(this.Ctx).Model(&m_account.Transfer{}).Where(filters).Find(&dbModels)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("transfers:fetch_failed", "获取交易失败"))
	}
	transfers := make([]*Transfer, 0)
	for _, dbModel := range dbModels{
		transfers = append(transfers, NewTransferFromModel(this.Ctx, dbModel))
	}
	return transfers
}

// GetFeeTransfersByThirdBids 根据订单号获取手续费交易记录
func (this *TransferRepository) GetFeeTransfersByThirdBids(thirdBIds []string) []*Transfer{
	feeAccount := b_account.NewAccountRepository(this.Ctx).GetFeeAccount()
	filters := eel.Map{
		"third_bid__in": thirdBIds,
		"dest_account_id": feeAccount.Id,
		"is_deleted": false,
	}
	return this.GetByFilters(filters)
}

func (this *TransferRepository) GetById(tid int) *Transfer{
	filters := eel.Map{
		"id": tid,
	}
	transfers := this.GetByFilters(filters)
	if len(transfers) > 0{
		return transfers[0]
	}
	return nil
}

func (this *TransferRepository) GetByThirdBid(bid string) []*Transfer{
	filters := eel.Map{
		"third_bid": bid,
	}
	return this.GetByFilters(filters)
}

func NewTransferRepository(ctx context.Context) *TransferRepository{
	instance := new(TransferRepository)
	instance.Ctx = ctx
	return instance
}