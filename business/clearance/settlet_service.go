package clearance

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/business/account"
	b_transfer "github.com/gingerxman/ginger-finance/business/transfer"
	m_clearance "github.com/gingerxman/ginger-finance/models/clearance"
	"github.com/gingerxman/gorm"
	// "github.com/go-redsync/redsync"
)

type SettleService struct{
	eel.ServiceBase
}

func NewSettleService(ctx context.Context) *SettleService{
	instance := new(SettleService)
	instance.Ctx = ctx
	return instance
}

func (this *SettleService) isOrderSettled(bid string) bool {
	return b_transfer.NewTransferRepository(this.Ctx).IsSettled(bid)
}

func (this *SettleService) DoSettle(recordIds []int) {
	if len(recordIds) == 0 {
		return
	}
	
	records := NewClearanceRepository(this.Ctx).GetRecordsByIds(recordIds)
	bid := records[0].Bid
	eel.Logger.Debug(bid)
	eel.Logger.Debug(this.isOrderSettled(bid))
	
	if this.isOrderSettled(bid) {
		return
	}
	
	if true {
		// TODO: 按业务编号加锁
		//lockKey := fmt.Sprintf("order_settlement_%s", bid)
		//lockOption := vanilla.NewLockOption(lockKey)
		//lockOption.SetTimeout(120)
		//mutex, err := vanilla.Lock.Lock(lockKey, lockOption)
		//defer func(m *redsync.Mutex) {
		//	if m != nil {
		//		m.Unlock()
		//	}
		//}(mutex)
		//if err != nil{
		//	panic(eel.NewBusinessError("redis_lock:get_failed", "获取redis锁失败"))
		//}
	}

	allAccountIds := make([]int, 0)
	for _, record := range records{
		allAccountIds = append(allAccountIds, record.SourceAccountId)
		allAccountIds = append(allAccountIds, record.DestAccountId)
	}
	id2account := account.NewAccountRepository(this.Ctx).GetId2Account(allAccountIds)

	params := make([]b_transfer.TransferParams, 0)
	
	for _, record := range records{
		businessData := record.GetBusinessData()
		trigger := "settle"
		if businessData.BizCode == "deposit" {
			trigger = "deposit"
		}
		action := fmt.Sprintf("%s: bid_%s", trigger, record.Bid)
		
		destAccount := id2account[record.DestAccountId]
		sourceAccount := id2account[record.SourceAccountId]
		
		params = append(params, b_transfer.TransferParams{
			SourceAccount: sourceAccount,
			DestAccount: destAccount,
			SourceAmount: record.Amount,
			DestAmount: record.Amount,
			Bid: record.Bid,
			Action: action,
			ExtraData: map[string]interface{}{
				"trigger": trigger,
				"source_imoney": record.SourceImoneyCode,
				"dest_imoney": record.DestImoneyCode,
				// "labels": record.ExtraData.Labels,
			},
		})

		//if destAccount != nil && mapset.NewSetFromSlice([]interface{}{"artist", "artist_manager"}).Contains(userType){
		//	notifingRecords = append(notifingRecords, accountParams.IncomeNotifyParams{
		//		UserId: destAccount.GetUserId(),
		//		UserType: userType,
		//		Bid: bid,
		//		Amount: record.Amount,
		//		Item: item,
		//	})
		//}
	}
	ll := len(params)
	if ll > 0{
		b_transfer.NewTransferService(this.Ctx).BulkTransfer(params)
		this.markRecordSettled(bid)
	}

	//if !reclear && len(notifingRecords) > 0{
	//	// 重新清算不发通知
	//	b_account.NewAccountBalanceService(this.Ctx).NotifyIncome(notifingRecords)
	//}

	// 结算完成消息
	//event.AsyncEvent.Send(events.SETTLEMENT_DONE, map[string]interface{}{
	//	"bid": bid,
	//})
}

// MarkClearedWithBid 标记清算记录已被结算
func (this *SettleService) markRecordSettled(bid string){
	db := eel.GetOrmFromContext(this.Ctx).Model(&m_clearance.ClearanceRecord{}).Where("bid", bid).Update(gorm.Params{
		"is_settled": true,
		"settle_status": m_clearance.SETTLE_STATUS_FINISHED,
	})
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("update_clearance_record_failed", "更新清算记录失败"))
	}
}
