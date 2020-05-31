package account

import (
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/gingerxman/ginger-finance/business"
	
	"github.com/gingerxman/eel"

	m_account "github.com/gingerxman/ginger-finance/models/account"
)

var PLATFORM_CASH_ACCOUNT *Account

type AccountRepository struct{
	eel.ServiceBase
}

func (this *AccountRepository) GetByFilters(filters map[string]interface{}) []*Account{
	o := eel.GetOrmFromContext(this.Ctx)
	var dbModels []*m_account.Account
	db := o.Model(&m_account.Account{}).Where(filters).Order("id desc").Find(&dbModels)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("account:fetch_failed", "查询账户失败"))
	}
	
	var accounts []*Account
	for _, dbModel := range dbModels{
		accounts = append(accounts, NewAccountFromModel(this.Ctx, dbModel))
	}
	return accounts
}

func (this *AccountRepository) GetByCode(accountCode string) *Account{
	return NewAccountFactory(this.Ctx).GetOrCreate(CreateAccountParams{
		AccountCode: accountCode,
	})
}

func (this *AccountRepository) GetById(accountId int) *Account{
	filters := map[string]interface{}{
		"id": accountId,
	}
	accounts := this.GetByFilters(filters)
	if len(accounts) > 0{
		return accounts[0]
	}else{
		return nil
	}
}

func (this *AccountRepository) GetByIds(accountIds []int) []*Account{
	if len(accountIds) == 0{
		return []*Account{}
	}
	filters := map[string]interface{}{
		"id__in": accountIds,
	}
	return this.GetByFilters(filters)
}

func (this *AccountRepository) GetByUser(user business.IUser, imoneyCode string) *Account{
	return this.GetByUserId(user.GetId(), imoneyCode)
}

func (this *AccountRepository) GetByUserId(userId int, imoneyCode string) *Account{
	var accountCode string
	accountCode = fmt.Sprintf("%s_%d", imoneyCode, userId)
	return this.GetByCode(accountCode)
}

func (this *AccountRepository) formatAccountCodesFromUserIds(userIds []int, imoneyCode string) []string{
	accountCodes := make([]string, 0)
	for _, userId := range userIds{
		accountCode := ""
		if userId == 0{
			accountCode = imoneyCode
		}else{
			accountCode = fmt.Sprintf("%s_%d", imoneyCode, userId)
		}

		accountCodes = append(accountCodes, accountCode)
	}
	return accountCodes
}

func (this *AccountRepository) GetByUserIds(userIds []int, imoneyCode string) []*Account{
	accountCodes := this.formatAccountCodesFromUserIds(userIds, imoneyCode)
	accounts := make([]*Account, 0)

	if len(accountCodes) == 0{
		return accounts
	}
	filters := map[string]interface{}{
		"code__in": accountCodes,
	}
	accounts = this.GetByFilters(filters)

	// 对于没有查找到的userId，进行创建
	accountFactory := NewAccountFactory(this.Ctx)
	fetchedUserIds := make([]interface{}, 0)
	for _, account := range accounts{
		fetchedUserIds = append(fetchedUserIds, account.UserId)
	}
	iUserIds := make([]interface{}, 0)
	for _, userId := range userIds{
		iUserIds = append(iUserIds, userId)
	}
	for _, userId := range (mapset.NewSetFromSlice(iUserIds).Difference(mapset.NewSetFromSlice(fetchedUserIds))).ToSlice(){
		iUserId := userId.(int)
		accounts = append(accounts, accountFactory.GetOrCreate(CreateAccountParams{
			UserId: iUserId,
			ImoneyCode: imoneyCode,
			AccountCode: fmt.Sprintf("%s_%d", imoneyCode, iUserId),
		}))
	}
	return accounts
}

// GetFeeAccount 获取手续费账户
func (this *AccountRepository) GetFeeAccount() *Account{
	return this.GetByCode("cash_fee")
}

func (this *AccountRepository) GetPlatformCashAccount() *Account{
	return PLATFORM_CASH_ACCOUNT
}

func (this *AccountRepository) GetId2Account(ids []int) map[int]*Account{
	filters := map[string]interface{}{
		"id__in": ids,
	}
	accounts := this.GetByFilters(filters)
	id2Account := make(map[int]*Account)
	for _, account := range accounts{
		id2Account[account.Id] = account
	}
	return id2Account
}

func NewAccountRepository(ctx context.Context) *AccountRepository{
	instance := new(AccountRepository)
	instance.Ctx = ctx
	return instance
}