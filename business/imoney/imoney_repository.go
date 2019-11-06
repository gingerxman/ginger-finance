package imoney

import (
	"context"
	
	"github.com/gingerxman/eel"

	m_imoney "github.com/gingerxman/ginger-finance/models/imoney"
)


type ImoneyRepository struct{
	eel.ServiceBase
}

func (this *ImoneyRepository) GetByFilters(filters map[string]interface{}) []*Imoney{
	o := eel.GetOrmFromContext(this.Ctx)
	var dbModels []*m_imoney.IMoney
	db := o.Model(&m_imoney.IMoney{}).Where(filters).Find(&dbModels)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("imoney:fetch_failed", "查询虚拟货币失败"))
	}
	imoneys := make([]*Imoney, 0, len(dbModels))
	for _, dbModel := range dbModels{
		imoneys = append(imoneys, NewImoneyFromModel(this.Ctx, dbModel))
	}
	return imoneys
}

func (this *ImoneyRepository) GetByCode(code string) *Imoney{
	imoneys := this.GetByFilters(eel.Map{
		"code": code,
	})
	if len(imoneys) > 0{
		return imoneys[0]
	}
	return nil
}

func (this *ImoneyRepository) GetFromRAM(code string) *Imoney{
	data := Code2Imoney[code]
	return NewImoneyFromMap(this.Ctx, data)
}

func NewImoneyRepository(ctx context.Context) *ImoneyRepository{
	instance := new(ImoneyRepository)
	instance.Ctx = ctx
	return instance
}