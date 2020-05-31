package imoney

import (
	"context"

	"github.com/gingerxman/eel"
)

type ImoneyRepository struct{
	eel.ServiceBase
	imoneyManager *ImoneyManager
}

func (this *ImoneyRepository) GetByCode(code string) *Imoney{
	return this.imoneyManager.GetImoneyByCode(code)
}

func NewImoneyRepository(ctx context.Context) *ImoneyRepository{
	instance := new(ImoneyRepository)
	instance.Ctx = ctx
	instance.imoneyManager = NewImoneyManager(ctx)
	return instance
}