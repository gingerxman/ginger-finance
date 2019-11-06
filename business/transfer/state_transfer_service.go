package transfer

import (
	"context"
	
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
)


type StateTransferService struct{
	eel.ServiceBase
}

// GetTotalExpenseForAccount 获取账户的总支出
func (this *StateTransferService) GetTotalExpenseForAccount(account *b_account.Account) int {
	value := 0
	sql := "SELECT SUM(source_amount) AS total_expense FROM account_transfer WHERE source_account_id=? AND is_deleted=0"
	sqlParams := []interface{}{account.Id}
	db := eel.GetOrmFromContext(this.Ctx).Raw(sql, sqlParams...).Scan(&value)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("state_expense:failed", "统计总支出失败"))
	}

	return value
}

// getTotalIncome
// beego的RawSet构建params方法有问题，导致group by查询只返回一条数据，这里不用params，直接将数据拼接到sql中
func (this *StateTransferService) getTotalIncome(filters eel.Map) map[int]int{
	return make(map[int]int, 0)
	/*
	sql := "SELECT dest_account_id, SUM(dest_amount) AS total_income FROM account_transfer WHERE is_deleted=0"

	if filters != nil && len(filters) > 0{
		if v, ok := filters["explicit_time"]; ok && v.(string) != ""{
			explicitTime := v.(string)
			sql += fmt.Sprintf(" AND created_at<='%s'", explicitTime)
		}
		if v, ok := filters["created_at__range"]; ok && len(v.([]string)) > 0{
			timeRange := v.([]string)
			sql += fmt.Sprintf(" AND created_at BETWEEN '%s' AND '%s'", timeRange[0], timeRange[1])
		}
		if v, ok := filters["account_id__in"]; ok && len(v.([]int)) > 0{
			ss := make([]string, 0)
			for _, iv := range v.([]int){
				ss = append(ss, strconv.Itoa(iv))
			}
			sql += fmt.Sprintf(" AND dest_account_id in (%s)", strings.Join(ss, ","))
		}
	}
	sql += " GROUP BY dest_account_id"
	var maps []orm.Params
	_, err := eel.GetOrmFromContext(this.Ctx).Raw(sql).Values(&maps)

	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("state_income:failed", "统计收入失败"))
	}

	accountId2totalIncome := make(map[int]decimal.Decimal)
	if len(maps) > 0{
		for _, data := range maps{
			if data["dest_account_id"] == nil{
				continue
			}
			accountId, _ := strconv.Atoi(data["dest_account_id"].(string))
			income := constant.DECIMAL_ZERO
			if v, ok := data["total_income"]; ok && v != nil{
				income, _ = decimal.NewFromString(v.(string))
			}
			accountId2totalIncome[accountId] = income
		}
	}
	return accountId2totalIncome
	 */
}

func (this *StateTransferService) GetTotalIncomeForAccounts(accountIds []int, filters eel.Map) map[int]int{
	if filters == nil{
		filters = eel.Map{}
	}
	if len(accountIds) == 0{
		return make(map[int]int)
	}

	filters["account_id__in"] = accountIds
	return this.getTotalIncome(filters)
}

// GetTotalIncomeForAccount 获取账户总收入
func (this *StateTransferService) GetTotalIncomeForAccount(accountId int, args ...map[string]interface{}) int{

	filters := eel.Map{
		"account_id__in": []int{accountId},
	}

	switch len(args) {
	case 1:
		iFilters := args[0]
		for k, v := range iFilters{
			filters[k] = v
		}
	}

	accountId2Income := this.getTotalIncome(filters)

	income := 0

	if v, ok := accountId2Income[accountId]; ok{
		income = v
	}

	return income
}


// GetRealtimeValidBalance
// 实时计算(特定时间的)可用余额 = 特定时间transfer总收入 - transfer总支出 - frozen总金额
func (this *StateTransferService) GetRealtimeValidBalance(account *b_account.Account, args ...map[string]interface{}) int{
	return 0
	/*
	explicitTime := ""
	switch len(args) {
	case 1:
		if v, ok := args[0]["explicit_time"]; ok && v.(string) != ""{
			explicitTime = v.(string)
		}
	}

	// 使用触发器情况下，并且没有时间限定，则直接返回
	//if explicitTime == "" && account.UseTrigger(){
	//	return account.Balance.Sub(account.FrozenAmount)
	//}

	totalIncome := this.GetTotalIncomeForAccount(account.Id, map[string]interface{}{
		"explicit_time": explicitTime,
	})
	totalExpense := this.GetTotalExpenseForAccount(account)

	balance := totalIncome.Sub(totalExpense)
	totalFrozenAmount := b_account.NewAccountBalanceService(this.Ctx).GetTotalFrozenAmountForAccount(account)

	// 更新余额
	if explicitTime == "" && !account.UseTrigger(){
		b_account.NewAccountBalanceService(this.Ctx).UpdateBalanceForAccount(account, balance, totalFrozenAmount)
	}

	beego.Info(fmt.Sprintf("account(%d)==> total_income: %s; total_expense: %s; total_frozen: %s",
		account.Id, totalIncome.String(), totalExpense.String(), totalFrozenAmount.String()))

	return balance.Sub(totalFrozenAmount).Round(2)
	 */
}

func NewStateTransferService(ctx context.Context) *StateTransferService{
	instance := new(StateTransferService)
	instance.Ctx = ctx
	return instance
}