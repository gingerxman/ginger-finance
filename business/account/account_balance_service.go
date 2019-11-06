package account

import (
	"context"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/business/account/params"
	m_account "github.com/gingerxman/ginger-finance/models/account"
)


type AccountBalanceService struct{
	eel.ServiceBase
}

//// UpdateBalanceForAccount 更新账户余额
//func (this *AccountBalanceService) UpdateBalanceForAccount(account *Account, newBalance, newFrozenAmount decimal.Decimal){
//	if account.UseTrigger(){
//		return
//	}
//	// 先加锁
//	qs := eel.GetOrmFromContext(this.Ctx).Model(&m_account.Account{}).Where("id", account.Id)
//	_ = qs.ForUpdate().One(&m_account.Account{})
//
//	oldBalance := account.Balance
//
//	newBalanceAmount, _ := newBalance.Round(2).Float64()
//	newFAmount, _ := newFrozenAmount.Round(2).Float64()
//	_, err := qs.Update(gorm.Params{
//		"balance": newBalanceAmount,
//		"frozen_amount": newFAmount,
//		"updated_at": time.Now(),
//	})
//
//	if err != nil{
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("account:update_balance_failed", "更新余额失败"))
//	}
//
//	if oldBalance.Equal(newBalance){
//		return
//	}
//	account.Balance = newBalance
//	this.balanceUpdated(account, oldBalance, newBalance)
//}
//
//// IncreaseBalanceForAccount 增加余额
//func (this *AccountBalanceService) IncreaseBalanceForAccount(account *Account, amount decimal.Decimal, iBusiness iBusiness.IBusiness){
//
//	if amount.LessThan(decimal.Zero){
//		panic(eel.NewBusinessError("increase_balance:failed", "数值不能为负"))
//	}
//
//	if account.ForgetAboutBalance(){
//		return
//	}
//
//	if !account.UseTrigger(){
//		_, err := eel.GetOrmFromContext(this.Ctx).Model(&m_account.Account{}).Where("id", account.Id).Update(gorm.Params{
//			"balance": orm.ColFloatValue(orm.ColAdd, amount),
//			"updated_at": time.Now(),
//		})
//		if err != nil{
//			eel.Logger.Error(err)
//			panic(eel.NewBusinessError("increase_balance:failed", "增加账户余额失败"))
//		}
//	}
//
//	oldBalance := account.Balance
//	account.Balance = account.Balance.Add(amount)
//	newBalance := account.Balance
//	this.balanceUpdated(account, oldBalance, newBalance, iBusiness)
//}
//
//// DecreaseBalanceForAccount 减少余额
//func (this *AccountBalanceService) DecreaseBalanceForAccount(account *Account, amount decimal.Decimal, iBusiness iBusiness.IBusiness){
//
//	if amount.LessThan(decimal.Zero){
//		panic(eel.NewBusinessError("increase_balance:failed", "数值不能为负"))
//	}
//
//	if account.ForgetAboutBalance(){
//		return
//	}
//
//	if !account.UseTrigger(){
//		fAmount, _ := amount.Float64()
//		// 扣减余额要考虑到账户是否可透支，针对不可透支的账户进行余额控制
//		updateSql := "UPDATE account_account SET balance=balance-?, updated_at=? WHERE id=?"
//		updateParams := []interface{}{fAmount, time.Now(), account.Id}
//		if iBusiness.GetActionType() != "re_settlement" && !account.CanOverdraw(){ // 重新清算行为允许余额为负
//			updateSql += " AND balance-frozen_amount>=?"
//			updateParams = append(updateParams, fAmount)
//		}
//		res, err := eel.GetOrmFromContext(this.Ctx).Raw(updateSql, updateParams...).Exec()
//
//		if err != nil{
//			eel.Logger.Error(err)
//			panic(eel.NewBusinessError("increase_balance:failed", "减少账户余额失败"))
//		}
//
//		updatedCount, _ := res.RowsAffected()
//		if updatedCount == 0{
//			panic(eel.NewBusinessError("increase_balance:failed", "账户余额不足"))
//		}
//	}
//
//	oldBalance := account.Balance
//	account.Balance = account.Balance.Sub(amount)
//	newBalance := account.Balance
//	this.balanceUpdated(account, oldBalance, newBalance, iBusiness)
//}

// 收入通知
func (this *AccountBalanceService) NotifyIncome(records []params.IncomeNotifyParams){
	eel.Logger.Warn("delete NotifyIncome()")
	//if common.NewEnvironService(this.Ctx).IsLocalEnv(){
	//	return
	//}
	//resource := eel.NewResource(this.Ctx)
	//for _, record := range records{
	//	_, err := resource.Put("weser", "template_message.income_message", map[string]interface{}{
	//		"user_id": record.UserId,
	//		"mp_account_name": beego.AppConfig.String("system::MP_ACCOUNT_NAME"),
	//		"order_bid": record.Bid,
	//		"money": record.Amount.Round(2).String(),
	//		"income_type": common.NewNameService(this.Ctx).GetAppDisplayName(record.Item),
	//	})
	//	if err != nil{
	//		eel.Logger.Error(err)
	//	}
	//}
}

// GetTotalFrozenAmountForAccount 获取账户冻结资产总额
func (this *AccountBalanceService) GetTotalFrozenAmountForAccount(account *Account) int {
	sql := "SELECT SUM(amount) AS total_amount FROM account_frozen_record WHERE account_id=? AND status=?"
	sqlParams := []interface{}{account.Id, m_account.ACCOUNT_FROZEN_STATUS["FROZEN"]}
	
	amount := 0
	db := eel.GetOrmFromContext(this.Ctx).Raw(sql, sqlParams...).Scan(&amount)
	err := db.Error

	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("state_frozen_amount:failed", "统计冻结资产失败"))
	}

	return amount
}

func NewAccountBalanceService(ctx context.Context) *AccountBalanceService{
	instance := new(AccountBalanceService)
	instance.Ctx = ctx
	return instance
}