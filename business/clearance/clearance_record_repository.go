package clearance

import (
	"context"
	
	"github.com/gingerxman/eel"
	m_clearance "github.com/gingerxman/ginger-finance/models/clearance"
)


type ClearanceRepository struct{
	eel.ServiceBase
}

func (this *ClearanceRepository) GetPagedRecordsByFilters(filters map[string]interface{}, pageInfo *eel.PageInfo) ([]*ClearanceRecord, eel.INextPageInfo){
	var dbModels []*m_clearance.ClearanceRecord
	qs := eel.GetOrmFromContext(this.Ctx).Model(&m_clearance.ClearanceRecord{}).Where(filters).Order("id desc")

	records := make([]*ClearanceRecord, 0)
	paginateResult, db := eel.Paginate(qs, pageInfo, &dbModels)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return records, paginateResult
	}

	for _, dbModel := range dbModels{
		records = append(records, NewClearanceRecordFromDbModel(this.Ctx, dbModel))
	}
	return records, paginateResult
}

func (this *ClearanceRepository) GetByFilters(filters eel.Map) []*ClearanceRecord{
	var dbModels []*m_clearance.ClearanceRecord
	db := eel.GetOrmFromContext(this.Ctx).Model(&m_clearance.ClearanceRecord{}).Where(filters).Order("id").Limit(-1).Find(&dbModels)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("clearance_record:fetch_failed", "获取结算记录失败"))
	}
	records := make([]*ClearanceRecord, 0)
	for _, dbModel := range dbModels{
		records = append(records, NewClearanceRecordFromDbModel(this.Ctx, dbModel))
	}
	return records
}

func (this *ClearanceRepository) GetByBid(bid string) []*ClearanceRecord{
	filters := eel.Map{
		"bid": bid,
	}
	return this.GetByFilters(filters)
}

// GetUnclearedRecords 为结算记录
func (this *ClearanceRepository) GetUnclearedRecords() []*ClearanceRecord{
	var dbModels []*m_clearance.ClearanceRecord
	filters := map[string]interface{}{
		"cleared": false,
	}
	db := eel.GetOrmFromContext(this.Ctx).Model(&m_clearance.ClearanceRecord{}).Where(filters).Order("id").Limit(-1).Find(&dbModels)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("clearance_record:get_failed", "获取未结算记录失败"))
	}

	records := make([]*ClearanceRecord, 0, len(dbModels))
	for _, dbModel := range dbModels{
		records = append(records, NewClearanceRecordFromDbModel(this.Ctx, dbModel))
	}
	return records
}

func (this *ClearanceRepository) GetRecordsByIds(ids []int) []*ClearanceRecord {
	filters := eel.Map{
		"id__in": ids,
	}
	
	return this.GetByFilters(filters)
}

func NewClearanceRepository(ctx context.Context) *ClearanceRepository{
	instance := new(ClearanceRepository)
	instance.Ctx = ctx
	return instance
}