package frozen

import (
	"context"
	"fmt"
	
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	"github.com/gingerxman/gorm"
	"strings"
	"time"
	
	b_transfer "github.com/gingerxman/ginger-finance/business/transfer"
	m_account "github.com/gingerxman/ginger-finance/models/account"
)


type FrozenAccountService struct{
	eel.ServiceBase
}

//// increaseFrozenAmountForAccount 增加冻结数额
//func (this *FrozenAccountService) increaseFrozenAmountForAccount(account *b_account.Account, amount int){
//	account.FrozenAmount = account.FrozenAmount + amount
//	if account.FrozenAmount > account.Balance && !account.CanOverdraw() {
//		panic(eel.NewBusinessError("not_enough_balance", "余额不足,无法完成资金冻结"))
//	}
//
//	db := eel.GetOrmFromContext(this.Ctx).Model(&m_account.Account{}).Where("id", account.Id).Update(gorm.Params{
//		"frozen_amount": gorm.Expr("frozen_amount + ?", amount),
//		"updated_at": time.Now(),
//	})
//	err := db.Error
//	if err != nil{
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("increase_frozen_amount:failed", "增加账户冻结金额失败"))
//	}
//}
//
//// releaseFrozenRecordForAccount 释放冻结资产
//func (this *FrozenAccountService) releaseFrozenRecordForAccount(account *b_account.Account, frozenRecord *FrozenRecord, toStatus int8){
//	o := eel.GetOrmFromContext(this.Ctx)
//	db := o.Model(&m_account.FrozenRecord{}).Where("id", frozenRecord.Id).Update(gorm.Params{
//		"status": toStatus,
//		"updated_at": time.Now(),
//	})
//	err := db.Error
//	frozenRecord.Status = toStatus
//
//	if err != nil{
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("unfrozen_account:failed", "解冻资产失败"))
//	}
//
//	//if account.UseTrigger(){
//	//	return
//	//}
//	//amount := frozenRecord.Amount
//	//account.FrozenAmount = account.FrozenAmount.Sub(amount)
//	//if account.FrozenAmount.LessThan(constant.DECIMAL_ZERO){
//	//	panic(eel.NewBusinessError("negative_frozen_amount", "冻结金额已成负值"))
//	//}
//	//
//	//fAmount, _ := amount.Float64()
//	//
//	//updateSql := "UPDATE account_account SET frozen_amount=frozen_amount-?, updated_at=? WHERE id=? AND frozen_amount>=?"
//	//updateParams := []interface{}{fAmount, time.Now(), account.Id, fAmount}
//	//
//	//res, err := o.Raw(updateSql, updateParams...).Exec()
//	//
//	//if err != nil{
//	//	eel.Logger.Error(err)
//	//	panic(eel.NewBusinessError("decrease_frozen_amount:failed", "减少账户冻结金额失败"))
//	//}
//	//
//	//updatedCount, _ := res.RowsAffected()
//	//if updatedCount == 0{
//	//	panic(eel.NewBusinessError("decrease_frozen_amount:failed", "账户冻结金额不足"))
//	//}
//}

// FrozenForAccount 冻结资产
func (this *FrozenAccountService) FrozenForAccount(account *b_account.Account, amount int, frozenType string, remark string) *FrozenRecord {
	dbModel := &m_account.FrozenRecord{
		AccountId:  account.Id,
		ImoneyCode: account.GetImoneyCode(),
		Amount:     amount,
		Type:       m_account.STR2FROZENTYPE[frozenType],
		Status:     m_account.ACCOUNT_FROZEN_STATUS["FROZEN"],
		ExtraData:  eel.ToJsonString(eel.Map{}),
	}
	db := eel.GetOrmFromContext(this.Ctx).Create(dbModel)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		errCode := "frozen_account:failed"
		if strings.Contains(err.Error(), "[mysql trigger]"){
			errCode = "frozen_record:not_enough_balance"
		}
		panic(eel.NewBusinessError(errCode, "保存冻结记录失败"))
	}
	return NewFrozenRecordFromModel(this.Ctx, dbModel)
}

// UnfrozenForAccount 解冻资产
func (this *FrozenAccountService) UnfrozenForAccount(account *b_account.Account, frozenRecord *FrozenRecord){
	toStatus := m_account.ACCOUNT_FROZEN_STATUS["UNFROZEN"]
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_account.FrozenRecord{}).Where("id", frozenRecord.Id).Update(gorm.Params{
		"status": toStatus,
		"updated_at": time.Now(),
	})
	
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("unfrozen_account:failed", "解冻资产失败"))
	}
	frozenRecord.Status = toStatus
}

func (this *FrozenAccountService) getDestAccount(frozenRecord *FrozenRecord) *b_account.Account {
	var destAccount *b_account.Account
	accountRepository := b_account.NewAccountRepository(this.Ctx)
	switch frozenRecord.FrozenType {
	case m_account.STR2FROZENTYPE["withdraw"]:
		destAccount = accountRepository.GetByCode(fmt.Sprintf("withdraw_%s", frozenRecord.ImoneyCode))
	case m_account.STR2FROZENTYPE["deduction"]:
		destAccount = accountRepository.GetByCode(frozenRecord.ImoneyCode)
	case m_account.STR2FROZENTYPE["consume"]:
		destAccount = accountRepository.GetByCode(frozenRecord.ImoneyCode)
	default:
		panic(eel.NewBusinessError("settleFrozenAccount:failed", "不合法的冻结记录"))
	}
	return destAccount
}

// SettlefrozenForAccount 消费冻结资产
func (this *FrozenAccountService) SettleFrozenForAccount(account *b_account.Account, frozenRecord *FrozenRecord, extraData map[string]interface{}) *b_transfer.Transfer{

	amount := frozenRecord.Amount
	destAccount := this.getDestAccount(frozenRecord)

	transferParams := b_transfer.TransferParams{
		SourceAccount: account,
		DestAccount: destAccount,
		SourceAmount: amount,
		DestAmount: amount,
		Bid: extraData["bid"].(string),
		Action: extraData["action"].(string),
		ExtraData: extraData,
	}
	transfer := b_transfer.NewTransferService(this.Ctx).Transfer(transferParams)

	o := eel.GetOrmFromContext(this.Ctx)
	var transferDbModel m_account.Transfer
	db := o.Model(&m_account.Transfer{}).Where("Bid", transfer.Bid).Take(&transferDbModel)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("unfrozen_account:failed", "查询交易失败"))
	}
	transfer.Id = transferDbModel.Id

	db = o.Model(&m_account.FrozenRecord{}).Where("id", frozenRecord.Id).Update(gorm.Params{
		"transfer_id": transfer.Id,
		"updated_at": time.Now(),
	})
	err = db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("unfrozen_account:failed", "解冻资产失败"))
	}
	frozenRecord.TransferId = transfer.Id

	return transfer
}

// RollbackSettledFrozenRecord 回滚到冻结状态
// 物理删除交易记录，恢复各账户余额
func (this *FrozenAccountService) RollbackSettledFrozenRecord(frozenRecord *FrozenRecord, transfer *b_transfer.Transfer){
	o := eel.GetOrmFromContext(this.Ctx)
	if frozenRecord.TransferId != 0 && transfer != nil && frozenRecord.TransferId == transfer.Id{
		db := o.Where("Id", frozenRecord.TransferId).Delete(&m_account.Transfer{})
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("rollback_frozen:failed", "回滚冻结消费失败"))
		}

	}

	frozenRecord.Status = m_account.ACCOUNT_FROZEN_STATUS["FROZEN"]
	frozenRecord.TransferId = 0
	frozenRecord.UpdatedAt = time.Now()

	db := o.Model(&m_account.FrozenRecord{}).Where("Id", frozenRecord.Id).Update(gorm.Params{
		"status": frozenRecord.Status,
		"transfer_id": 0,
		"updated_at": frozenRecord.UpdatedAt,
	})
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("rollback_frozen:failed", "回滚冻结消费失败"))
	}
}

func NewFrozenAccountService(ctx context.Context) *FrozenAccountService{
	instance := new(FrozenAccountService)
	instance.Ctx = ctx
	return instance
}