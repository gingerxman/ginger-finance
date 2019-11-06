package frozen

import (
	"context"
	
	"github.com/gingerxman/eel"

	m_account "github.com/gingerxman/ginger-finance/models/account"
)

type FrozenRecordRepository struct{
	eel.ServiceBase
}

func (this *FrozenRecordRepository) GetByFilters(filters map[string]interface{}) []*FrozenRecord {
	o := eel.GetOrmFromContext(this.Ctx)
	var dbModels []*m_account.FrozenRecord
	db := o.Model(&m_account.FrozenRecord{}).Where(filters).Find(&dbModels)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("account_frozen_record:fetch_failed", "获取资产冻结记录失败"))
	}
	var records []*FrozenRecord
	for _, dbModel := range dbModels{
		records = append(records, NewFrozenRecordFromModel(this.Ctx, dbModel))
	}
	return records
}

func (this *FrozenRecordRepository) GetById(recordId int) *FrozenRecord {
	filters := map[string]interface{}{
		"id": recordId,
	}
	records := this.GetByFilters(filters)
	if len(records) > 0{
		return records[0]
	}else{
		return nil
	}
}

func NewFrozenRecordRepository(ctx context.Context) *FrozenRecordRepository{
	instance := new(FrozenRecordRepository)
	instance.Ctx = ctx
	return instance
}