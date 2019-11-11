package params

import (
	string_util "github.com/gingerxman/ginger-finance/business/common/util"
	"reflect"
	"strings"
)

type ClearanceRecordParams struct{
	Name string `json:"name"`
	Bid string	`json:"bid"`
	Amount int	`json:"amount"`
	Ratio float64 `json:"ratio"`
	SourceImoneyCode string `json:"source_imoney_code,omitempty"`
	ImoneyCode string `json:"imoney_code"`
	SourceUserId int	`json:"source_user_id,omitempty"`
	DestUserId int	`json:"dest_user_id,omitempty"`
	DestUserRole string `json:"dest_user_role,omitempty"`

	OrderParams *OrderParams `json:"order_params,omitempty"`

	Action string `json:"action,omitempty"`
	SourceAccountId int `json:"source_account_id,omitempty"`
	DestAccountId int `json:"dest_account_id,omitempty"`
}

func (this *ClearanceRecordParams) GetFieldValueByKey(field string) interface{}{
	field = string_util.ToCamelString(field)
	return reflect.Indirect(reflect.ValueOf(this)).FieldByName(field).Interface()
}

func (this *ClearanceRecordParams) GetRuleName() string{
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

func NewClearanceRecordParamsFromClearanceParams(clearancePrams *ClearanceParams) *ClearanceRecordParams{
	params :=  &ClearanceRecordParams{
		Name:             clearancePrams.Name,
		Bid:              clearancePrams.Bid,
		Amount:           clearancePrams.Amount,
		Ratio:            1.0,
		SourceImoneyCode: clearancePrams.SourceImoneyCode,
		ImoneyCode:       clearancePrams.ImoneyCode,
		SourceUserId:     clearancePrams.SourceUserId,
		DestUserId:       clearancePrams.DestUserId,
		DestUserRole:     "member",
		OrderParams:      clearancePrams.OrderParams,
		Action:           "",
		SourceAccountId:  0,
		DestAccountId:    0,
	}
	if params.SourceImoneyCode == ""{
		params.SourceImoneyCode = params.ImoneyCode
	}
	return params
}