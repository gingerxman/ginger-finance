package clearance

import (
	string_util "github.com/gingerxman/ginger-finance/business/common/util"
	"reflect"
	"strconv"
	"strings"
)

type ClearanceData struct{
	Bid string	`json:"bid"`
	Amount int	`json:"amount"`
	ImoneyCode string `json:"imoney_code"`
	Name string	`json:"name"`
	SourceUserId int	`json:"source_user_id,omitempty"`
	DestUserId int	`json:"dest_user_id,omitempty"`
	Ratio float64 `json:"ratio,omitempty"`

	FeeUserId int `json:"fee_user_id,omitempty"`

	SourceAccountId int `json:"source_account_id,omitempty"`
	DestAccountId int `json:"dest_account_id,omitempty"`
}

// IsSysRule 是否系统内置的规则
func (this *ClearanceData) IsSysRule() bool{
	splits := strings.Split(this.Name, ".")
	return splits[0] == "sys"
}

func (this *ClearanceData) IsDepositData() bool{
	return this.Name == "sys.deposit" || this.Name == "sys.mpcoin_product"
}

func (this *ClearanceData) NeedSyncSettlement() bool{
	switch this.Name {
	case "sys.bomb", "lucky_money":
		return true
	}
	return false
}

// GetTargetUserId 获取目标用户id
// 目标用户是指当前的清算数据要使用的规则的所有者
func (this *ClearanceData) GetTargetUserId() int{
	splits := strings.Split(this.Name, ".")
	if splits[0] == "sys"{
		return 0
	}else{
		userId, _ := strconv.Atoi(splits[0])
		return userId
	}
}

func (this *ClearanceData) GetRuleName() string{
	return strings.Split(this.Name, ".")[1]
}

func (this *ClearanceData) GetFieldValueByKey(field string) interface{}{
	field = string_util.ToCamelString(field)
	return reflect.Indirect(reflect.ValueOf(this)).FieldByName(field).Interface()
}