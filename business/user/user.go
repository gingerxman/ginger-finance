package user

import (
	"context"
	"github.com/gingerxman/eel"
	"time"

	"github.com/bitly/go-simplejson"
	
)

type User struct {
	eel.EntityBase
	Id                int
	PlatformId        int
	Unionid           string
	Name              string
	Avatar            string
	Sex               string
	Code              string
	RawData           *simplejson.Json
	CreatedAt         time.Time
}

func NewUserFromOnlyId(ctx context.Context, id int) *User {
	user := new(User)
	user.Ctx = ctx
	user.Model = nil
	user.Id = id
	return user
}

func GetUserFromContext(ctx context.Context) *User {
	user := ctx.Value("user").(*User)
	return user
}

func (this *User) GetId() int {
	return this.Id
}

func (this *User) GetName() string {
	return this.Name
}

func (this *User) GetAvatar() string {
	return this.Avatar
}

func init() {
}
