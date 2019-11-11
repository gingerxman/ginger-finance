package clearance

import (
	"context"
	"encoding/json"
	
	"github.com/gingerxman/eel"
	"time"
	
	b_account "github.com/gingerxman/ginger-finance/business/account"
	m_clearance "github.com/gingerxman/ginger-finance/models/clearance"
)

type businessData struct {
	BizCode string `json:biz_code`
	BizData map[string]interface{} `json:biz_data`
}

type ClearanceRecord struct {
	eel.EntityBase

	Id int
	SourceUserId int
	SourceAccountId int
	SourceImoneyCode string
	DestUserId int
	DestAccountId int
	DestImoneyCode string
	Amount int
	Ratio float64
	Bid string
	extraData string
	IsSettled bool
	CreatedAt time.Time
	SettledAt time.Time
	
	businessData *businessData

	// need fill
	SourceAccount *b_account.Account
	DestAccount *b_account.Account
}

func (this *ClearanceRecord) GetRawExtraData() string{
	return this.extraData
}

func (this *ClearanceRecord) GetBusinessData() *businessData {
	if this.extraData == ""{
		return nil
	}
	data := businessData{}
	err := json.Unmarshal([]byte(this.extraData), &data)
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("clearance_record:decode_business_data_failed", "解析BusinessData失败"))
	}
	return &data
}

func NewClearanceRecordFromDbModel(ctx context.Context, dbModel *m_clearance.ClearanceRecord) *ClearanceRecord{
	instance := new(ClearanceRecord)
	instance.Id = dbModel.Id
	instance.Bid = dbModel.Bid
	instance.SourceUserId = dbModel.SourceUserId
	instance.SourceAccountId = dbModel.SourceAccountId
	instance.SourceImoneyCode = dbModel.SourceImoneyCode
	instance.DestUserId = dbModel.DestUserId
	instance.DestAccountId = dbModel.DestAccountId
	instance.DestImoneyCode = dbModel.DestImoneyCode
	instance.Ratio = dbModel.Ratio
	instance.Amount = dbModel.Amount
	instance.IsSettled = dbModel.IsSettled
	// instance.UserType = dbModel.UserType
	instance.extraData = dbModel.ExtraData
	instance.CreatedAt = dbModel.CreatedAt
	instance.SettledAt = dbModel.SettledAt

	return instance
}