package imoney

import (
	"github.com/gingerxman/eel"
	"time"
)

// IMoney 虚拟资产
type IMoney struct {
	eel.Model
	Code string `gorm:"size:52"`
	DisplayName string `gorm:"size:52"`
	ExchangeRate float64 `gorm:"default:1.0"`
	IsPayable bool `gorm:"default:true"`
	IsDebtable bool `gorm:"default:false"`
	IsEnabled bool `gorm:"default:false"`
}
func (this *IMoney) TableName() string {
	return "imoney_imoney"
}


var IMONEY_WITHDRAW_STATUS = map[string]int8{
	"REQUEST": 1,
	"TRANSFERRING": 2,
	"SUCCESS": 3,
	"FAILED": 4,
	"REJECTED": 5,
}
var WITHDRAW_STATUS2STR = map[int8]string{
	IMONEY_WITHDRAW_STATUS["REQUEST"]: "request",
	IMONEY_WITHDRAW_STATUS["SUCCESS"]: "success",
	IMONEY_WITHDRAW_STATUS["REJECTED"]: "rejected",
	IMONEY_WITHDRAW_STATUS["FAILED"]: "failed",
	IMONEY_WITHDRAW_STATUS["TRANSFERRING"]: "transferring",
}
var WITHDRAW_STR2STATUS = map[string]int8{
	"request": IMONEY_WITHDRAW_STATUS["REQUEST"],
	"success": IMONEY_WITHDRAW_STATUS["SUCCESS"],
	"reject": IMONEY_WITHDRAW_STATUS["REJECTED"],
	"failed": IMONEY_WITHDRAW_STATUS["FAILED"],
	"transferring": IMONEY_WITHDRAW_STATUS["TRANSFERRING"],
}
var WITHDRAW_CHANNEL = map[string]int8{
	"UNKNOW": -1,
	"WEIXIN": 0,
	"BANK": 1,
}
var WITHDRAW_CAHNNEL2TEXT = map[int8]string{
	-1: "unknow",
	0: "weixin",
	1: "bank",
}
type Withdraw struct{
	eel.Model
	AccountId int
	Bid string `gorm:"size:125;unique"`
	Status int8 `gorm:"default:1"`
	ImoneyCode string `gorm:"size:52;index"`
	Amount int
	WithdrawRate float64 `gorm:"type:decimal(5,2)"`
	Money int
	Channel int8 `gorm:"default:0"`
	OuterTransferId string `gorm:"size:125"` // 第三方交易号
	TransferId int
	FrozenRecordId int
	Remark string `gorm:"size:256"`
	ExtraData string `gorm:"type:text"`
	FinishedAt time.Time `gorm:"type:datetime"`
}
func (this *Withdraw) TableName() string {
	return "imoney_withdraw"
}





func init() {
	eel.RegisterModel(new(IMoney))
	eel.RegisterModel(new(Withdraw))
}
