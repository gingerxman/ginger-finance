package account

import (
	"context"
	"fmt"
	"github.com/gingerxman/ginger-finance/business"
	"strconv"
	"strings"

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
			params.AccountCode = fmt.Sprintf("%s.user.%d", imoneyCode, userId)
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
		if err := db.Error; err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("account:get_failed", "获取财务账户失败"))
		}
	}else{
		var isDebtable bool
		var imoneyCode string
		if params.WithDebtable{
			isDebtable = params.IsDebtable
		}else{
			if params.ImoneyCode == ""{
				imoneyCode = strings.Split(params.AccountCode, ".")[0]
			}else{
				imoneyCode = params.ImoneyCode
			}
			// 平台账户都可负债
			if imoneyCode == params.AccountCode{
				isDebtable = true
			}else{
				imoney := b_imoney.NewImoneyRepository(this.Ctx).GetByCode(imoneyCode)
				if imoney == nil {
					panic(eel.NewBusinessError("imoney:invalid_imoney", "不合法的虚拟资产: "+ imoneyCode))
				}
				isDebtable = imoney.IsDebtable
			}
		}
		
		if params.UserId == 0 && imoneyCode != params.AccountCode {
			sps := strings.Split(params.AccountCode, ".")
			uid, err := strconv.Atoi(sps[len(sps) - 1])
			if err != nil{
				eel.Logger.Error(err)
				panic(eel.NewBusinessError("account:parse_user_id_failed", "解析用户id失败"))
			}
			params.UserId = uid
		}

		dbModel = m_account.Account{
			Code: params.AccountCode,
			UserId: params.UserId,
			ImoneyCode: imoneyCode,
			IsDebtable: isDebtable,
			IsDeleted: false,
		}
		db := o.Create(&dbModel)
		if err := db.Error; err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("account:create_failed", "创建财务账户失败"))
		}
	}
	return NewAccountFromModel(this.Ctx, &dbModel)
}

func NewAccountFactory(ctx context.Context) *AccountFactory{
	instance := new(AccountFactory)
	instance.Ctx = ctx
	return instance
}


func init(){
}