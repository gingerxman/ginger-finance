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

const SETTLE_STATUS_WAIT = 0
const SETTLE_STATUS_RUN = 1
const SETTLE_STATUS_FINISHED = 2
var status2str = map[int]string {
	SETTLE_STATUS_WAIT: "wait",
	SETTLE_STATUS_RUN: "run",
	SETTLE_STATUS_FINISHED: "wait",
}
var str2status = map[string]int {
	"wait": SETTLE_STATUS_WAIT,
	"run": SETTLE_STATUS_RUN,
	"finished": SETTLE_STATUS_FINISHED,
}
// ClearanceRecord 清算记录
type ClearanceRecord struct {
	eel.Model
	SourceUserId int `gorm:"index"`
	SourceAccountId int `gorm:"index"`
	SourceImoneyCode string
	DestUserId int `gorm:"index"`
	DestAccountId int `gorm:"index"`
	DestImoneyCode string
	Amount int
	Ratio float64 `gorm:"type:decimal(5,2);default:1.0"`
	Bid string `gorm:"size:125;unique"`
	//UserType string `gorm:"size:128);default(unknown)"`
	Remark string `gorm:"size:256"`
	ExtraData string `gorm:"type:text"`
	IsSettled bool `gorm:"index"`
	SettleStatus int `gorm:"index"`
	SettledAt time.Time `gorm:"type:datetime"`
}
func (this *ClearanceRecord) TableName() string {
	return "clearing_record"
}




func init() {
	eel.RegisterModel(new(ClearanceRule))
	eel.RegisterModel(new(ClearanceRecord))
}
