package clearance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	b_user "github.com/gingerxman/ginger-finance/business/user"
	
	"github.com/gingerxman/eel"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	"github.com/gingerxman/ginger-finance/business/clearance/params"
	m_clearance "github.com/gingerxman/ginger-finance/models/clearance"
)

//type clearanceRuleEngine interface {
//	GetName() string
//	DoClearance(context.Context, *params.ClearanceParams) []*PreparedClearanceRecord
//}
//
//var name2ClearanceRuleEngine = make(map[string]clearanceRuleEngine)
//func RegisterClearanceRuleEngine(engine clearanceRuleEngine){
//	name := engine.GetName()
//	name2ClearanceRuleEngine[name] = engine
//}

type clearanceEngine struct{
	eel.ServiceBase

	params *params.ClearanceParams

	// feeAmount int // 手续费
	
	accountRepository *b_account.AccountRepository
}

// getOrderParams 从ginger-mall service获取订单信息
func (this *clearanceEngine) getOrderParams(orderBid, orderStatus string) *params.OrderParams{
	resp, err := eel.NewResource(this.Ctx).Get("ginger-mall", "order.order", eel.Map{
		"bid": orderBid,
	})
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order:get_order_failed", "获取订单详情失败"))
	}

	orderData := resp.Data().MustMap()
	if orderData != nil{
		if orderStatus == "all" || orderData["status"].(string) == orderStatus{
			orderParams := new(params.OrderParams)
			orderBytes, _ := json.Marshal(orderData)
			err := json.Unmarshal(orderBytes, orderParams)
			if err != nil{
				eel.Logger.Error(err)
				panic(eel.NewBusinessError("order_params:decode_failed", "解析订单数据失败"))
			}
			if orderParams.CorpId != 0{
				corp := b_user.NewCorpFromOnlyId(this.Ctx, orderParams.CorpId)
				orderParams.CorpRelatedUserId = corp.GetRelatedUser().Id
			}
			return orderParams
		} else {
			panic(eel.NewBusinessError("order:not_exist", fmt.Sprintf("订单不存在 (status:%s, bid:%s)", orderStatus, orderBid)))
		}
	} else {
		panic(eel.NewBusinessError("order:not_exist", fmt.Sprintf("订单不存在 (status:%s, bid:%s)", orderStatus, orderBid)))
	}
}

// ParseOrder 根据订单信息解析得到ClearanceParams数据
func (this *clearanceEngine) ParseOrder(orderBid, orderStatus string, args ...string) *params.ClearanceParams{
	orderParams := this.getOrderParams(orderBid, orderStatus)
	
	productName := orderParams.Invoices[0].Products[0].Name
	engineName := "order"
	
	clearanceParams := &params.ClearanceParams{
		Name: fmt.Sprintf("%s.%s", engineName, productName),
		Bid: orderBid,
		OrderParams: orderParams,
		Amount: orderParams.FinalMoney,
		ImoneyCode: "rmb",
		SourceUserId: orderParams.UserId, // 默认清算金额来源为下单用户
	}
	
	this.params = clearanceParams
	return clearanceParams
}

// DoOrderDefaultClearance 订单默认清算
// orderUser.rmb => supplier.cash
func (this *clearanceEngine) doOrderClearance() *ClearanceRecord {
	// 首先扣减手续费
	orderParams := this.params.OrderParams

	amount := orderParams.FinalMoney
	if amount <= 0 {
		return nil
	}
	
	sourceImoneyCode := "rmb"
	destImoneyCode := "cash"
	sourceAccount := this.accountRepository.GetByUserId(orderParams.UserId, sourceImoneyCode)
	destAccount := this.accountRepository.GetByUserId(orderParams.CorpRelatedUserId, destImoneyCode)

	ratio := 1.0
	return &ClearanceRecord{
		SourceUserId: orderParams.UserId,
		SourceAccountId: sourceAccount.Id,
		SourceImoneyCode: sourceImoneyCode,
		DestUserId: orderParams.CorpRelatedUserId,
		DestAccountId: destAccount.Id,
		DestImoneyCode: destImoneyCode,
		Amount: amount,
		Ratio: ratio,
		Bid: orderParams.Bid,
		businessData: &businessData{
			BizCode: "order",
			BizData: map[string]interface{} {
				"bid": orderParams.Bid,
				"final_money": orderParams.FinalMoney,
			},
		},
	}
}

// saveRecords 存储清算记录
func (this *clearanceEngine) saveRecords(records []*ClearanceRecord) []int {
	spew.Dump(records)
	models := make([]*m_clearance.ClearanceRecord, 0, len(records))
	for _, record := range records{
		if record.SourceAccountId == record.DestAccountId{
			continue
		}
		if record.Amount <= 0{
			continue
		}
		models = append(models, &m_clearance.ClearanceRecord{
			Bid: record.Bid,
			SourceUserId: record.SourceUserId,
			SourceAccountId: record.SourceAccountId,
			SourceImoneyCode: record.SourceImoneyCode,
			DestUserId: record.DestUserId,
			DestAccountId: record.DestAccountId,
			DestImoneyCode: record.DestImoneyCode,
			Amount: record.Amount,
			Ratio: record.Ratio,
			IsSettled: false,
			SettleStatus: m_clearance.SETTLE_STATUS_WAIT,
			ExtraData: eel.ToJsonString(record.businessData),
		})
	}

	ids := make([]int, 0)
	for _, model := range models {
		db := eel.GetOrmFromContext(this.Ctx).Create(model)
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("clearance:failed", "保存清算数据失败"))
		}
		ids = append(ids, model.Id)
	}
	
	return ids
}

// OrderIsCleared 检查订单是否已清算
func (this *clearanceEngine) IsOrderCleared(orderBid string) bool{
	return eel.GetOrmFromContext(this.Ctx).Model(&m_clearance.ClearanceRecord{}).Where("bid", orderBid).Exist()
}

// DoClearance
func (this *clearanceEngine) DoClearance(orderBid, orderStatus string, args ...string) []int {
	records := make([]*ClearanceRecord, 0)

	if this.params == nil{
		this.ParseOrder(orderBid, orderStatus, args...)
	}
	clearanceParams := this.params
	if clearanceParams == nil{
		return make([]int, 0)
	}
	// 加锁
	// TODO: 改进LockOption
	//if clearanceParams.SourceUserId != 0{
	//	lockKey := fmt.Sprintf("order_clearance_user_%d", clearanceParams.SourceUserId)
	//	lockOption := eel.NewLockOption(lockKey)
	//	lockOption.SetTimeout(60)
	//	mutex, err := vanilla.Lock.Lock(lockKey, lockOption)
	//	defer func(m *redsync.Mutex) {
	//		if m != nil {
	//			m.Unlock()
	//		}
	//	}(mutex)
	//	if err != nil{
	//		panic(eel.NewBusinessError("redis_lock:get_failed", "获取redis锁失败"))
	//	}
	//
	//}

	clearanceRecord := this.doOrderClearance()
	if clearanceRecord != nil {
		records = append(records, clearanceRecord)
	}
	return this.saveRecords(records)
}

func NewClearanceEngine(ctx context.Context) *clearanceEngine{
	instance := new(clearanceEngine)
	instance.Ctx = ctx
	instance.accountRepository = b_account.NewAccountRepository(ctx)
	return instance
}