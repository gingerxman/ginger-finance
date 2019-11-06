package user

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodeUserService struct {
	eel.ServiceBase
}

func NewEncodeUserService(ctx context.Context) *EncodeUserService {
	service := new(EncodeUserService)
	service.Ctx = ctx
	return service
}

func (this *EncodeUserService) Encode(user *User) *RUser {
	return &RUser{
		Id:                user.Id,
		Name:              user.Name,
		Avatar:            user.Avatar,
		Sex:               user.Sex,
		Code:              user.Code,
	}
}

func (this *EncodeUserService) EncodeMany(users []*User) []*RUser {
	rows := make([]*RUser, len(users))
	
	for i, user := range users {
		rows[i] = this.Encode(user)
	}
	
	return rows
}
