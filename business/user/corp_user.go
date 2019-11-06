package user

import (
	"context"
	"github.com/gingerxman/eel"
	"time"
)

type CorpUser struct {
	eel.EntityBase
	Id         int
	CreatedAt  time.Time
}

func (this *CorpUser) GetId() int {
	return this.Id
}

func NewCorpUserFromId(ctx context.Context, corpUserId int) *CorpUser{
	instance := new(CorpUser)
	instance.Ctx = ctx
	instance.Id = corpUserId
	return instance
}

func init() {
}
