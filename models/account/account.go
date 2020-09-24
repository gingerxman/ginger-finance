package account

import (
	"github.com/gingerxman/eel"
	"time"
)

// Account 财务账户
type Account struct {
	eel.Model
	Code string `gorm:"size:52;unique"`
	UserId int `gorm:"index"`
	ImoneyCode string `gorm:"size:52"`
	AccountType string `gorm:"size:64"` // 账户类型，如corp、sys、supplier ...
	Balance int `gorm:"default:0"` // 余额
	FrozenAmount int `gorm:"default:0"` // 冻结数额
	IsDebtable bool `gorm:"default:false"`	// 是否可负债
	IsDeleted bool `gorm:"default:false"`
}
func (this *Account) TableName() string {
	return "account_account"
}

// AccountBalanceChangeLog 余额变动日志
type AccountBalanceChangeLog struct{
	Id int
	SourceAccountId int `gorm:"index"`
	DestAccountId int `gorm:"index"`
	TransferId int `gorm:"index"`
	Action string `gorm:"size(128);default('')"` // 交易类别，比如 withdraw, transfer等
	ThirdBid string `gorm:"size(256);default('')"` // 业务编号
	SourceAccountBalance float64 `gorm:"default(0.00);digits(12);decimals(2)"`
	DestAccountBalance float64 `gorm:"default(0.00);digits(12);decimals(2)"`
	CreatedAt time.Time `gorm:"auto_now_add;type(datetime)"`
}
func (this *AccountBalanceChangeLog) TableName() string{
	return "account_balance_change_log"
}

const FR_STATUS_FROZEN = 1
const FR_STATUS_UNFROZEN = 2
const FR_STATUS_SETTLED = 3
var ACCOUNT_FROZEN_STATUS = map[string]int8{
	"FROZEN": FR_STATUS_FROZEN, // 冻结
	"UNFROZEN": FR_STATUS_UNFROZEN, // 已解冻
	"SETTLED": FR_STATUS_SETTLED, // 已消费
}

var STR2FROZENTYPE = map[string]int8{
	"unknown": 0, // 未知类型
	"consume": 1, // 订单虚拟资产
	"withdraw": 2, // 提现
	"deduction": 3, // 扣款
}
var FROZENTYPE2STR = map[int8]string {
	0: "unknown",
	1: "consume",
	2: "withdraw",
	3: "dedution",
}

// FrozenRecord 冻结记录
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
	ExtraData string `gorm:"type:text"` // 自描述形式的内容
	Digest string `gorm:"size:125;unique"`
	IsDeleted bool
}
func (this *Transfer) TableName() string {
	return "account_transfer"
}




type SystemTriggerLog struct{
	Id int `gorm:"primary_key"`
	RowId int
	Code string `orm:"size(125);default('')"`
	Name string `orm:"size(1024);default('')"`
	Status string `orm:"size(128);default('')"`
	Msg string `orm:"size(1024);default('')"`
	CreatedAt  time.Time `orm:"auto_now_add;type(datetime)"`
}
func (this *SystemTriggerLog) TableName() string{
	return "system_trigger_log"
}

func init() {
	eel.RegisterModel(new(Account))
	eel.RegisterModel(new(AccountBalanceChangeLog))
	eel.RegisterModel(new(FrozenRecord))
	eel.RegisterModel(new(Transfer))
	eel.RegisterModel(new(SystemTriggerLog))
}
