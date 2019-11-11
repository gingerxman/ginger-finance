package account

import (
	"context"
	"github.com/gingerxman/eel"
	m_account "github.com/gingerxman/ginger-finance/models/account"
)

type Account struct {
	eel.EntityBase

	Id int
	Code string
	IsDebtable bool
	Balance int
	FrozenAmount int
}

// GetImoneyCode 获取账户的虚拟货币
// 目前账户的Code有以下几种形式，可以从中解析出imoney_code
// 1. user_cash_xx
// 2. cash
// 3. cash.xx
// 4. cash.coupon.xx
func (this *Account) GetImoneyCode() string{
	return NewParseAccountCodeService(this.Ctx).ParseImoneyCodeFromCode(this.Code)
}

func (this *Account) GetUserId() int{
	//return this.Model.(*m_account.Account).UserId
	return NewParseAccountCodeService(this.Ctx).ParseUserIdFromCode(this.Code)
}

// IsPlatformNormalAccount 是否平台一般账户
func (this *Account) IsPlatformNormalAccount() bool{
	return this.GetUserId() == 0
}

// CanOverdraw 是否允许透支
// 平台一般账户都可以透支
func (this *Account) CanOverdraw() bool{
	if this.IsPlatformNormalAccount(){
		return true
	}
	return this.IsDebtable
}

// ForgetAboutBalance 是否可以不用关心余额
func (this *Account) ForgetAboutBalance() bool{
	return this.IsPlatformNormalAccount() || this.GetImoneyCode() == "rmb"
}

func (this *Account) GetValidBalance() int {
	return this.Balance - this.FrozenAmount
}

// GetFreshBalance 从数据库中获取余额
// 当前mysql的隔离级别为read-committed，为最大限度的避免幻读(并不能完全避免)，
// 在使用account对象的balance属性时，要重新从数据库中获取
//func (this *Account) GetFreshBalance() decimal.Decimal{
//	o := eel.GetOrmFromContext(this.Ctx)
//	var dbModel m_account.Account
//	err := o.Model(&m_account.Account{}).Where("Id", this.Id).One(&dbModel)
//	if err != nil{
//		eel.Logger.Error(err)
//		panic(eel.NewBusinessError("account:get_failed", "获取账户信息失败"))
//	}
//	return decimal.NewFromFloat(dbModel.Balance)
//}

func (this *Account) GetBalance() int {
	return this.Balance
}

func (this *Account) GetFrozenAmount() int {
	return this.FrozenAmount
}

func NewAccountFromModel(ctx context.Context, dbModel *m_account.Account) *Account{
	instance := new(Account)
	instance.Ctx = ctx
	instance.Model = dbModel

	instance.Id = dbModel.Id
	instance.Code = dbModel.Code
	instance.IsDebtable = dbModel.IsDebtable
	instance.Balance = dbModel.Balance - dbModel.FrozenAmount
	instance.FrozenAmount = dbModel.FrozenAmount
	
	return instance
}