package user

import (
	"context"
	"github.com/gingerxman/eel"
)

const _PLATFORM_CORP_ID = -1

type Corp struct {
	eel.EntityBase
	Id int
	
	PlatformId int
	CorpUser *CorpUser
}

func GetCorpFromContext(ctx context.Context) *Corp {
	val := ctx.Value("corp")
	if val == nil {
		return nil
	}
	corp := val.(*Corp)
	return corp
}

func (this *Corp) GetId() int {
	return this.Id
}

func (this *Corp) IsValid() bool {
	return this.Id != 0
}


func (this *Corp) IsPlatform() bool {
	return this.Id == _PLATFORM_CORP_ID
	//return this.Id == this.PlatformId
}

func (this *Corp) GetPlatformId() int {
	return _PLATFORM_CORP_ID
	//return this.PlatformId
}

func (this *Corp) GetRelatedUser() *User {
	resp, err := eel.NewResource(this.Ctx).LoginAsManager().Get("ginger-account", "corp.related_user_id", eel.Map{
		"corp_id": this.Id,
	})
	
	if err != nil {
		eel.Logger.Error(err)
		return nil
	}
	
	respData := resp.Data()
	userId := respData.Get("id").MustInt()
	
	return NewUserFromOnlyId(this.Ctx, userId)
}

func GetPlatformId() int {
	return _PLATFORM_CORP_ID
	//return this.PlatformId
}

func NewCorpFromOnlyId(ctx context.Context, id int) *Corp {
	instance := new(Corp)
	instance.Ctx = ctx
	instance.Model = nil
	instance.Id = id
	return instance
}

func NewInvalidCorp(ctx context.Context) *Corp {
	return NewCorpFromOnlyId(ctx, 0)
}

func init() {
}
