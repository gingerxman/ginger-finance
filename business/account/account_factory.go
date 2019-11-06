package account

import (
	"context"
	"fmt"
	"github.com/gingerxman/ginger-finance/business"
	
	"github.com/gingerxman/eel"
	b_imoney "github.com/gingerxman/ginger-finance/business/imoney"
	m_account "github.com/gingerxman/ginger-finance/models/account"
)

type CreateAccountParams struct{
	UserId int
	ImoneyCode string
	AccountCode string
	IsDebtable bool
	WithDebtable bool
}

var imoney2debtable = make(map[string]bool) // IMoney是否可负债

type AccountFactory struct{
	eel.ServiceBase
}

func (this *AccountFactory) CreateForUser(user business.IUser, params CreateAccountParams) *Account{

	accountCode := params.AccountCode
	userId := user.GetId()
	params.UserId = userId

	if accountCode == ""{
		if params.ImoneyCode == ""{
			panic(eel.NewBusinessError("account:create_failed", "创建财务账户失败,不合法的参数"))
		}
		imoneyCode := params.ImoneyCode
		if userId == 0{
			params.AccountCode = imoneyCode
		}else{
			params.AccountCode = fmt.Sprintf("%s_%d", imoneyCode, userId)
		}
	}

	return this.GetOrCreate(params)
}

func (this *AccountFactory) GetOrCreate(params CreateAccountParams) *Account{
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_account.Account{}).Where("code", params.AccountCode)
	
	var dbModel m_account.Account
	if db.Exist(){
		db := db.Take(&dbModel)
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("account:get_failed", "获取财务账户失败"))
		}
	}else{
		var isDebtable bool
		if params.WithDebtable{
			isDebtable = params.IsDebtable
		}else{
			// 如果当前的accountCode==imoneyCode，则为true, 否则以imoney的配置为准
			var imoneyCode string
			if params.ImoneyCode == ""{
				imoneyCode = NewParseAccountCodeService(this.Ctx).ParseImoneyCodeFromCode(params.AccountCode)
			}else{
				imoneyCode = params.ImoneyCode
			}

			if imoneyCode == params.AccountCode{
				isDebtable = true
			}else{
				imoney := b_imoney.NewImoneyRepository(this.Ctx).GetByCode(imoneyCode)
				if imoney == nil {
					panic(eel.NewBusinessError("imoney:invalid_imoney", imoneyCode))
				}
				
				isDebtable = imoney.IsDebtable
			}
		}
		
		if params.UserId == 0{
			params.UserId = NewParseAccountCodeService(this.Ctx).ParseUserIdFromCode(params.AccountCode)
		}

		dbModel = m_account.Account{
			UserId: params.UserId,
			Code: params.AccountCode,
			IsDeleted: false,
			IsDebtable: isDebtable,
		}
		db := o.Create(&dbModel)
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("account:create_failed", "创建财务账户失败"))
		}
	}
	return NewAccountFromModel(this.Ctx, &dbModel)
}

func loadIMoneys() {
	for _, imoney := range b_imoney.Code2Imoney{
		isDebtable := false
		if val, ok := imoney["is_debtable"]; ok{
			isDebtable = val.(bool)
		}
		imoney2debtable[imoney["code"].(string)] = isDebtable
	}
}


func NewAccountFactory(ctx context.Context) *AccountFactory{
	instance := new(AccountFactory)
	instance.Ctx = ctx
	return instance
}


func init(){
	// 初始化全局变量 imoney2debtable
	loadIMoneys()
}