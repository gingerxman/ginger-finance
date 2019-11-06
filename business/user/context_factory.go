package user

import (
	"context"
	"github.com/gingerxman/eel"
	
	"github.com/bitly/go-simplejson"
	
	"net/http"
)

func NewBusinessContext(ctx context.Context, request *http.Request, userId int, jwtToken string, rawData *simplejson.Json) context.Context {
	user := new(User)
	user.Model = nil
	user.Id = userId
	user.RawData = rawData
	
	ctx = context.WithValue(ctx, "jwt", jwtToken)
	user.Ctx = ctx
	
	ctx = context.WithValue(ctx, "user", user)
	
	//创建corp
	if rawData != nil {
		jwtType, _ := rawData.Get("type").Int()
		if jwtType == 1 { //for logined corp user
			corpId, err := rawData.Get("cid").Int()
			if err == nil {
				corp := new(Corp)
				corp.Model = nil
				corp.Id = corpId
				corp.Ctx = ctx
				
				//处理corp user
				corpUserId, _ := rawData.Get("uid").Int()
				corpUser := NewCorpUserFromId(ctx, corpUserId)
				corp.CorpUser = corpUser
				
				ctx = context.WithValue(ctx, "corp", corp)
			} else {
				eel.Logger.Error(err)
			}
		} else if jwtType == 2 { //for logined mall mobile user
			corpId, err := rawData.Get("cid").Int()
			if err == nil {
				corp := new(Corp)
				corp.Model = nil
				corp.Id = corpId
				corp.Ctx = ctx
				
				ctx = context.WithValue(ctx, "corp", corp)
			} else {
				eel.Logger.Error(err)
			}
		}
	}
	return ctx
}

func init() {
}
