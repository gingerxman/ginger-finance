package clearance

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/business/constant"
)


type EncodeClearanceRecordService struct{
	eel.ServiceBase
}

func (this *EncodeClearanceRecordService) Encode(record *ClearanceRecord) *RRecord{
	if record.DestAccount == nil{
		NewFillClearanceRecordService(this.Ctx).Fill([]*ClearanceRecord{record}, eel.FillOption{
			"with_account": true,
		})
	}
	return &RRecord{
		Id: record.Id,
		Bid: record.Bid,
		SourceUserId: record.SourceAccount.UserId,
		DestUserId: record.DestAccount.UserId,
		Role: "unknown",
		Amount: record.Amount,
		Ratio: record.Ratio,
		SettledAt: record.SettledAt.Format(constant.TIME_LAYOUT),
	}
}

func (this *EncodeClearanceRecordService) EncodeMany(records []*ClearanceRecord) []*RRecord{
	rRecords := make([]*RRecord, 0)
	if len(records) == 0{
		return rRecords
	}
	NewFillClearanceRecordService(this.Ctx).Fill(records, eel.FillOption{
		"with_account": true,
	})

	for _, record := range records{
		rRecords = append(rRecords, this.Encode(record))
	}
	return rRecords
}

//func (this *EncodeClearanceRecordService) EncodeManyGroupByBid(records []*ClearanceRecord) []map[string]interface{}{
//	groupedRecords := make([]map[string]interface{}, 0)
//	if len(records) == 0{
//		return groupedRecords
//	}
//	rRecords := this.EncodeMany(records)
//	bid2Records := make(map[string][]*RRecord)
//	bidHasFeeRecord := make(map[string]bool)
//	bids := make([]string, 0)
//	for _, rr := range rRecords{
//		bids = append(bids, rr.Bid)
//		if _, ok := bid2Records[rr.Bid]; ok{
//			bid2Records[rr.Bid] = append(bid2Records[rr.Bid], rr)
//		}else{
//			bid2Records[rr.Bid] = []*RRecord{rr}
//		}
//		if rr.Role == "sys_fee"{
//			bidHasFeeRecord[rr.Bid] = true
//		}
//	}
//	feeTransfers := b_transfer.NewTransferRepository(this.Ctx).GetFeeTransfersByThirdBids(bids)
//	bid2FeeTransfer := make(map[string]*b_transfer.Transfer)
//	for _, tr := range feeTransfers{
//		bid2FeeTransfer[tr.ThirdBid] = tr
//	}
//	for bid, records := range bid2Records{
//		if t, ok := bid2FeeTransfer[bid]; ok{
//			if v, ok2 := bidHasFeeRecord[bid]; !ok2 && !v{
//				records = append(records, &RRecord{
//					Bid: bid,
//					Role: "sys_fee",
//					Amount: t.DestAmount,
//					Ratio: 0.006,
//					SettledAt: t.CreatedAt.Format(constant.TIME_LAYOUT),
//				})
//			}
//		}
//		groupedRecords = append(groupedRecords, map[string]interface{}{
//			"bid": bid,
//			"settlements": records,
//		})
//	}
//	return groupedRecords
//}

func NewEncodeClearanceRecordService(ctx context.Context) *EncodeClearanceRecordService{
	instance := new(EncodeClearanceRecordService)
	instance.Ctx = ctx
	return instance
}