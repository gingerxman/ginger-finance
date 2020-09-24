package transfer

import (
	"context"
	"fmt"
	"github.com/gingerxman/ginger-finance/business/common/util"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/snowflake"
	b_account "github.com/gingerxman/ginger-finance/business/account"
	b_imoney "github.com/gingerxman/ginger-finance/business/imoney"
	m_account "github.com/gingerxman/ginger-finance/models/account"
	"sort"
	"strings"
)

var snowflakeNode, _ = snowflake.NewNode(1)

// TransferParams 交易的必要因素
type TransferParams struct{
	SourceAccount *b_account.Account
	DestAccount *b_account.Account
	SourceAmount int
	DestAmount int
	Bid string
	Action string
	ExtraData map[string]interface{}
}

type TransferService struct{
	eel.ServiceBase
}

// handlePreparedTransfers 批量存储、更新余额
func (this *TransferService) handlePreparedTransfers(preparedTransfers []*Transfer){
	createList := make([]*m_account.Transfer, 0)
	accountIds := make([]int, 0)
	accountId2account := make(map[int]*b_account.Account)
	for _, preparedTransfer := range preparedTransfers{
		createList = append(createList, &m_account.Transfer{
			Bid: preparedTransfer.Bid,
			ThirdBid: preparedTransfer.ThirdBid,
			SourceAccountId: preparedTransfer.SourceAccount.Id,
			DestAccountId: preparedTransfer.DestAccount.Id,
			SourceAmount: preparedTransfer.SourceAmount,
			DestAmount: preparedTransfer.DestAmount,
			Action: preparedTransfer.Action,
			Digest: preparedTransfer.Digest,
			IsDeleted: false,
			ExtraData: preparedTransfer.ExtraData,
		})
		accountIds = append(accountIds, preparedTransfer.SourceAccountId)
		accountIds = append(accountIds, preparedTransfer.DestAccountId)
		accountId2account[preparedTransfer.SourceAccountId] = preparedTransfer.SourceAccount
		accountId2account[preparedTransfer.DestAccountId] = preparedTransfer.DestAccount
	}
	l := len(createList)
	if l > 0{
		o := eel.GetOrmFromContext(this.Ctx)

		// 按序加锁
		// 考虑到触发器， 先加锁后插入
		sort.Ints(accountIds)
		tmp := map[int]bool{}
		for _, accountId := range accountIds{
			account := accountId2account[accountId]
			if account.ForgetAboutBalance(){
				continue
			}
			if _, ok := tmp[accountId]; !ok{
				tmp[accountId] = true
				db := o.Raw("select * from account_account where id = ? for update")
				err := db.Error
				if err != nil {
					eel.Logger.Error(err)
				}
			}
		}

		for _, item := range createList {
			db := o.Create(&item)
			err := db.Error
			if err != nil{
				eel.Logger.Error(err)
				panic(eel.NewBusinessError("transfer:save_failed", "存储交易失败"))
			}
		}
		//_, err := o.InsertMulti(l, createList)
		//if err != nil{
		//	eel.Logger.Error(err)
		//	panic(eel.NewBusinessError("transfer:save_failed", "存储交易失败"))
		//}

		//accountBalanceService := b_account.NewAccountBalanceService(this.Ctx)
		//for _, transfer := range preparedTransfers{
		//	accountBalanceService.DecreaseBalanceForAccount(transfer.SourceAccount, transfer.SourceAmount, transfer)
		//	accountBalanceService.IncreaseBalanceForAccount(transfer.DestAccount, transfer.DestAmount, transfer)
		//}
	}
}

// BulkTransfer 批量交易
func (this *TransferService) BulkTransfer(params []TransferParams) []*Transfer{
	preparedTransfers := make([]*Transfer, 0)
	for _, param := range params{
		preparedTransfers = append(preparedTransfers, this.prepareTransfer(param))
	}
	this.handlePreparedTransfers(preparedTransfers)
	return preparedTransfers
}

// Transfer 账户间交易
func (this *TransferService) Transfer(params TransferParams) *Transfer{
	preparedTransfer := this.prepareTransfer(params)
	this.handlePreparedTransfers([]*Transfer{preparedTransfer})

	return preparedTransfer
}

func (this *TransferService) generateBid() string {
	result := snowflakeNode.Generate().String()
	return result
}

func (this *TransferService) prepareTransfer(params TransferParams) *Transfer{
	sourceAccount := params.SourceAccount
	destAccount := params.DestAccount
	sourceAmount := params.SourceAmount
	destAmount := params.DestAmount
	thirdBid := params.Bid

	// 不同虚拟资产交易需要兑换
	sourceImoneyCode := sourceAccount.ImoneyCode
	destImoneyCode := destAccount.ImoneyCode
	if sourceAmount == destAmount && sourceImoneyCode != destImoneyCode{ // sourceAmount和destAmount不等时，认为已经兑换过
		destAmount = b_imoney.NewImoneyExchangeService(this.Ctx).ExchangeByImoneyCode(sourceImoneyCode, destImoneyCode, sourceAmount)
	}

	if thirdBid == ""{
		panic(eel.NewBusinessError("transfer:invalid_bid", "不合法的交易号"))
	}
	if sourceAccount == nil || destAccount == nil{
		panic(eel.NewBusinessError("transfer:invalid_account", "账户不存在"))
	}
	if sourceAmount < 0 || destAmount < 0 {
		panic(eel.NewBusinessError("transfer:invalid_amount", "不合法的交易金额"))
	}
	extraData := params.ExtraData
	action := params.Action
	if !strings.HasPrefix(action, "re_settlement") && !sourceAccount.CanOverdraw(){
		filters := map[string]interface{}{
			"id": sourceAccount.Id,
			"balance__gte": sourceAmount,
		}
		validAccounts := b_account.NewAccountRepository(this.Ctx).GetByFilters(filters)
		if len(validAccounts) == 0{
			panic(eel.NewBusinessError("transfer:not_enough_balance", "账户余额不足"))
		}
	}

	if action == ""{
		action = fmt.Sprintf("direct: bid_%s", thirdBid)
	}

	hashData := fmt.Sprintf("%d_%f_%d_%f_%s_%s", sourceAccount.Id, sourceAmount, destAccount.Id, destAmount, thirdBid, action)
	digest := util.GetMD5(hashData)

	return NewTransferFromMap(map[string]interface{}{
		"bid": this.generateBid(),
		"third_bid": thirdBid,
		"source_account": sourceAccount,
		"dest_account": destAccount,
		"source_amount": sourceAmount,
		"dest_amount": destAmount,
		"action": action,
		"extra_data": eel.ToJsonString(extraData),
		"digest": digest,
	})
}

// TransferForFee Deprecated
func (this *TransferService) TransferForFee(deductAccount *b_account.Account, feeAmount int, orderBid string) *Transfer{
	destAccount := b_account.NewAccountRepository(this.Ctx).GetFeeAccount()
	params := TransferParams{
		SourceAccount: deductAccount,
		DestAccount: destAccount,
		SourceAmount: feeAmount,
		DestAmount: feeAmount,
		Bid: orderBid,
		Action: fmt.Sprintf("settlement_fee: bid_%s", orderBid),
		ExtraData: map[string]interface{}{
			"code": destAccount.ImoneyCode,
			"user_type": "sys_fee",
			"item": "fee",
		},
	}
	return this.Transfer(params)
}

// TransferForDeposit 充值交易
// 1、user的rmb账户 => 平台的cash账户  订单结算中完成
// 2、平台的imoney账户 => 用户的imoney账户
func (this *TransferService) TransferForDeposit(destAccount *b_account.Account, amount int, orderBid string) *Transfer{
	imoneyCode := destAccount.ImoneyCode
	sourceAccount := b_account.NewAccountRepository(this.Ctx).GetByCode(imoneyCode)
	params := TransferParams{
		SourceAccount: sourceAccount,
		DestAccount: destAccount,
		SourceAmount: amount,
		DestAmount: amount,
		Bid: orderBid,
		Action: fmt.Sprintf("imoney.deposit: bid_%s", orderBid),
		ExtraData: map[string]interface{}{
			"code": imoneyCode,
			"user_type": "member",
			"item": "deposit",
		},
	}
	return this.Transfer(params)
}

func NewTransferService(ctx context.Context) *TransferService{
	instance := new(TransferService)
	instance.Ctx = ctx
	return instance
}