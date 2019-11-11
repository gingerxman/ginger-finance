package clearance

import (
	"encoding/json"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/business/clearance/params"
)

type PreparedClearanceRecordRuleNode struct{
	Role string	`json:"role"`
	Ratio float64	`json:"ratio"`
	AccountId int	`json:"account_id"`
}

type PreparedClearanceRecordExtraData struct {
	SettlementData *ClearanceData `json:"settlement_data,omitempty"`
	ClearanceData *params.ClearanceRecordParams `json:"clearance_data,omitempty"`
	Rule *params.FormattedRule `json:"rule,omitempty"`
	Node *PreparedClearanceRecordRuleNode `json:"node,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type PreparedClearanceRecord struct {
	SourceAccountId int
	DestAccountId int
	Amount int
	Ratio float64
	Bid string
	ExtraData *PreparedClearanceRecordExtraData
}

func NewPreparedClearanceRecordFromRecord(record *ClearanceRecord) *PreparedClearanceRecord{
	extraData := new(PreparedClearanceRecordExtraData)
	err := json.Unmarshal([]byte(record.GetRawExtraData()), extraData)
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("parse_clearance_extra_data:failed", "反序列化extra_data失败"))
	}
	return &PreparedClearanceRecord{
		SourceAccountId: record.SourceAccountId,
		DestAccountId: record.DestAccountId,
		Amount: record.Amount,
		Ratio: record.Ratio,
		Bid: record.Bid,
		ExtraData: extraData,
	}
}

func NewPreparedClearanceRecordFromRecordParams(recordParams *params.ClearanceRecordParams) *PreparedClearanceRecord{
	return &PreparedClearanceRecord{
		SourceAccountId: recordParams.SourceAccountId,
		DestAccountId: recordParams.DestAccountId,
		Bid: recordParams.Bid,
		Ratio: recordParams.Ratio,
		Amount: recordParams.Amount,
		ExtraData: &PreparedClearanceRecordExtraData{
			ClearanceData: recordParams,
			Node: &PreparedClearanceRecordRuleNode{
				Role: recordParams.DestUserRole,
				Ratio: recordParams.Ratio,
				AccountId: recordParams.DestAccountId,
			},
		},
	}
}