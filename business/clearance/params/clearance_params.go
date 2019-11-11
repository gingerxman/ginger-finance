package params

import (
	"strconv"
	"strings"
)

type ClearanceParams struct{
	Name string `json:"name"`
	Bid string	`json:"bid"`
	Amount int	`json:"amount"`
	SourceImoneyCode string `json:"source_imoney_code,omitempty"`
	ImoneyCode string `json:"imoney_code"`
	SourceUserId int	`json:"source_user_id,omitempty"`
	DestUserId int	`json:"dest_user_id,omitempty"`

	OrderParams *OrderParams `json:"order_params,omitempty"`
}

// IsSysRule 是否系统内置的规则
func (this *ClearanceParams) IsSysRule() bool{
	splits := strings.Split(this.Name, ".")
	return splits[0] == "sys"
}

// IsAppRule 是否互动规则
func (this *ClearanceParams) IsAppRule() bool{
	return !this.IsSysRule()
}

// IsAppLimitedRule 是否互动限制规则，只能适用默认订单清算
func (this *ClearanceParams) IsAppLimitedRule() bool{
	if strings.HasPrefix(this.Name, "sys"){
		return false
	}
	switch this.GetRuleName() {
	case "lucky_money", "lure", "member_vip":
		return true
	}

	return false
}

func (this *ClearanceParams) NeedSyncSettlement() bool{
	switch this.Name {
	case "sys.bomb", "lucky_money", "sys.imoney_deposit":
		return true
	}
	return false
}

// GetTargetUserId 获取目标用户id
// 目标用户是指当前的清算数据要使用的规则的所有者
func (this *ClearanceParams) GetTargetUserId() int{
	splits := strings.Split(this.Name, ".")
	if splits[0] == "sys"{
		return 0
	}else{
		userId, _ := strconv.Atoi(splits[0])
		return userId
	}
}

func (this *ClearanceParams) GetRuleName() string{
	sps := strings.Split(this.Name, ".")
	switch len(sps) {
	case 1:
		return sps[0]
	case 2:
		return sps[1]
	default:
		return this.Name
	}
}