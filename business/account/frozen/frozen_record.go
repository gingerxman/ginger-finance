package frozen

import (
	"context"
	"github.com/gingerxman/eel"
	m_account "github.com/gingerxman/ginger-finance/models/account"
	"time"
)

type FrozenRecord struct {
	eel.EntityBase

	Id int
	AccountId int
	ImoneyCode string
	Amount int
	Status int8
	TransferId int
	ExtraData string
	UpdatedAt time.Time
	CreatedAt time.Time
	FrozenType int8
}

// IsSettled 是否已消费
func (this *FrozenRecord) IsSettled() bool{
	return this.Status == m_account.ACCOUNT_FROZEN_STATUS["SETTLED"]
}

func NewFrozenRecordFromModel(ctx context.Context, dbModel *m_account.FrozenRecord) *FrozenRecord{
	instance := new(FrozenRecord)
	instance.Ctx = ctx
	instance.Model = dbModel

	instance.Id = dbModel.Id
	instance.AccountId = dbModel.AccountId
	instance.ImoneyCode = dbModel.ImoneyCode
	instance.Amount = dbModel.Amount
	instance.Status = dbModel.Status
	instance.TransferId = dbModel.TransferId
	instance.ExtraData = dbModel.ExtraData
	instance.UpdatedAt = dbModel.UpdatedAt
	instance.CreatedAt = dbModel.CreatedAt
	instance.FrozenType = dbModel.Type

	return instance
}