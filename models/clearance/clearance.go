package clearance

import (
	"github.com/gingerxman/eel"
	"time"
)

// ClearanceRule 清算规则
type ClearanceRule struct {
	eel.Model
	UserId int
	AccountId int
	Name string `gorm:"size:128;index"`
	ParentRuleId int
	Rule string `gorm:"type:text"`
}
func (this *ClearanceRule) TableName() string {
	return "clearing_rule"
}

// ClearanceRecord 清算记录
type ClearanceRecord struct {
	eel.Model
	SourceAccountId int
	DestAccountId int
	Amount int
	Ratio float64 `gorm:"type:decimal(5,2);default:1.0"`
	Bid string `gorm:"size:125;unique"`
	IsCleared bool
	//UserType string `gorm:"size:128);default(unknown)"`
	Remark string `gorm:"size:256"`
	ExtraData string `gorm:"type:text"`
	ClearedAt time.Time `gorm:"type:datetime"`
}
func (this *ClearanceRecord) TableName() string {
	return "clearing_record"
}




func init() {
	eel.RegisterModel(new(ClearanceRule))
	eel.RegisterModel(new(ClearanceRecord))
}
