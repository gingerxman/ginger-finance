package imoney

import (
	"context"
	"github.com/gingerxman/eel"
)


type EncodeImoneyService struct{
	eel.ServiceBase
}

func (this *EncodeImoneyService) Encode(imoney *Imoney) *RImoney{
	return &RImoney{
		Id: imoney.Id,
		Code: imoney.Code,
		DisplayName: imoney.DisplayName,
		ExchangeRate: imoney.ExchangeRate,
		EnableFraction: imoney.EnableFraction,
		IsPayable: imoney.IsPayable,
		IsDebtable: imoney.IsDebtable,
	}
}

func (this *EncodeImoneyService) EncodeMany(imoneys []*Imoney) []*RImoney{
	rImoneys := make([]*RImoney, 0)
	for _, imoney := range imoneys{
		rImoneys = append(rImoneys, this.Encode(imoney))
	}
	return rImoneys
}

func NewEncodeImoneyService(ctx context.Context) *EncodeImoneyService{
	instance := new(EncodeImoneyService)
	instance.Ctx = ctx
	return instance
}