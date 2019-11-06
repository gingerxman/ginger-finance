package product

import (
	"github.com/gingerxman/eel"
	"time"
)

// Account 财务账户
type Account struct {
	eel.Model
	Code string `gorm:"size:52;unique"`
	UserId int `gorm:"index"`
	Balance int // 余额
	FrozenAmount int // 冻结数额
	IsDebtable bool `gorm:"default:false"`	// 是否可负债
	IsDeleted bool `gorm:"default:false"`
	
	Path string `gorm:"size:256"`
}
func (this *Account) TableName() string {
	return "account_account"
}

// Transfer 交易
type Transfer struct{
	eel.Model
	Bid string `gorm:"size:256;"` // 交易编号
	ThirdBid string `gorm:"size:256"` // 如果该交易有关联的第三方交易，这里存储第三方交易业务编号
	SourceAccountId int `gorm:"index"`
	DestAccountId int `gorm:"index"`
	SourceAmount int
	DestAmount int
	Action string `gorm:"size:128"`
	Description string `gorm:"type:text"`
	Digest string `gorm:"size:125;unique"`
	IsDeleted bool
}
func (this *Transfer) TableName() string {
	return "account_transfer"
}


// FrozenRecord 冻结记录
const FR_STATUS_FROZEN = 1
const FR_STATUS_UNFROZEN = 2
const FR_STATUS_SETTLED = 3
var ACCOUNT_FROZEN_STATUS = map[string]int8{
	"FROZEN": FR_STATUS_FROZEN, // 冻结
	"UNFROZEN": FR_STATUS_UNFROZEN, // 已解冻
	"SETTLED": FR_STATUS_SETTLED, // 已消费
}
var FROZEN_TYPE = map[string]int8{
	"UNKNOWN": 0, // 未知类型
	"ORDER_IMONEY": 1, // 订单虚拟资产
	"WITHDRAW": 2, // 提现
	"DEDUCTION": 3, // 扣款
}
type FrozenRecord struct{
	eel.Model
	AccountId int `gorm:"index"`
	ImoneyCode string
	Amount int
	Type int8
	Status int8 `gorm:"default:1;index"`
	TransferId int
	ExtraData string `gorm:"type(text)"`
}
func (this *FrozenRecord) TableName() string {
	return "account_frozen_record"
}


type SystemTriggerLog struct{
	Id int `gorm:"primary_key"`
	RowId int
	Name string `orm:"size(1024);default('')"`
	Status string `orm:"size(128);default('')"`
	Msg string `orm:"size(1024);default('')"`
	CreatedAt  time.Time `orm:"auto_now_add;type(datetime)"`
}
func (this *SystemTriggerLog) TableName() string{
	return "system_trigger_log"
}


//type BalanceChangeLog struct{
//	Id int
//	AccountId int `gorm:"index"`
//	TransferId int `gorm:"index"`
//	BusinessType string `gorm:"size(128);default('')"` // 交易类别，比如 withdraw, transfer等
//	BusinessBid string `gorm:"size(256);default('')"` // 交易编号
//	SourceBalance float64 `gorm:"default(0.00);digits(12);decimals(2)"`
//	DestBalance float64 `gorm:"default(0.00);digits(12);decimals(2)"`
//	CreatedAt time.Time `gorm:"auto_now_add;type(datetime)"`
//}
//func (this *BalanceChangeLog) TableName() string{
//	return "account_balance_change_log"
//}





func init() {
	eel.RegisterModel(new(Account))
	eel.RegisterModel(new(Transfer))
	eel.RegisterModel(new(FrozenRecord))
	eel.RegisterModel(new(SystemTriggerLog))
}
