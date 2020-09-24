package account

import (
	"context"
	"github.com/gingerxman/eel"
	m_account "github.com/gingerxman/ginger-finance/models/account"
)

type Account struct {
	eel.EntityBase

	Id int
	UserId int
	AccountType string
	Code string
	ImoneyCode string
	IsDebtable bool
	Balance int
	FrozenAmount int
}

// IsSysNormalAccount 是否系统一般账户
func (this *Account) IsSysNormalAccount() bool{
	return this.UserId == 0
}

// CanOverdraw 是否允许透支
// 平台一般账户都可以透支
func (this *Account) CanOverdraw() bool{
	if this.IsSysNormalAccount(){
		return true
	}
	return this.IsDebtable
}

// ForgetAboutBalance 是否可以不用关心余额(即在交易后不更新余额)
func (this *Account) ForgetAboutBalance() bool{
	return this.IsSysNormalAccount()
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
	instance.UserId = dbModel.UserId
	instance.AccountType = dbModel.AccountType
	instance.Code = dbModel.Code
	instance.ImoneyCode = dbModel.ImoneyCode
	instance.IsDebtable = dbModel.IsDebtable
	instance.Balance = dbModel.Balance - dbModel.FrozenAmount
	instance.FrozenAmount = dbModel.FrozenAmount
	
	return instance
}