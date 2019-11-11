package params

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	
	"github.com/gingerxman/eel"
)

type OrderParams struct{
	Bid string `json:"bid"`
	CorpId int `json:"corp_id"`
	UserId int `json:"user_id"`
	FinalMoney int `json:"final_money"`
	Invoices []*invoice `json:"invoices"`
	ExtraData *orderExtraData `json:"extra_data"`

	CorpRelatedUserId int `json:"-"`
}

type invoice struct{
	Products []*product `json:"products"`
}

type product struct{
	Name string `json:"name"`
}

type orderExtraData struct{
	RelevantUserId int `json:"relevant_user_id"`
	RelevantCorpId int `json:"relevant_corp_id"`
	ClearanceRule *orderClearanceRule `json:"settlement_rule"`
	DepositImoneyData *imoneyDepositData `json:"deposit_imoney"`
	NoResettle bool `json:"no_resettle"`
}

type imoneyDepositData struct{
	Code string `json:"code"`
	Amount int `json:"amount"`
}

type orderClearanceRule struct{
	Name string `json:"name"`
	ImoneyCode string `json:"imoney_code"`
	Amount int `json:"amount"`
	SourceUserId int `json:"source_user_id"`
	DestUserId int `json:"dest_user_id"`
	Ratio float64 `json:"ratio"`
}

func NewOrderForClearanceParamsFromJson(order *simplejson.Json) *OrderParams{

	params :=  new(OrderParams)
	orderBytes, _ := order.MarshalJSON()
	err := json.Unmarshal(orderBytes, params)
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order_for_clearance_params:decode_failed", "解析订单数据失败"))
	}

	return params
}