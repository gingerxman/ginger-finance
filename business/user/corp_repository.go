package user

import (
	"context"
	"github.com/gingerxman/eel"
)

type CorpRepository struct {
	eel.ServiceBase
}

func NewCorpRepository(ctx context.Context) *CorpRepository {
	instance := new(CorpRepository)
	instance.Ctx = ctx
	return instance
}

func (this *CorpRepository) GetCorpById(corpId int) *Corp {
	return NewCorpFromOnlyId(this.Ctx, corpId)
}

//func (this *CorpRepository) GetCorpByUuid(uuid string) *Corp {
//	resp, err := eel.NewResource(this.Ctx).LoginAsManager().Get("gskep", "corp.corp", eel.Map{
//		"uuid": uuid,
//	})
//	if err != nil {
//		eel.Logger.Error(err)
//		return nil
//	} else {
//		respData := resp.Data()
//		corpData := respData.Get("corp")
//		return NewCorpFromOnlyId(this.Ctx, corpData.Get("id").MustInt())
//	}
//}
//
//var _PLATFORM_CORP map[string]interface{} = nil
//var _PLATFORM_CORP_ID int = 0
//
//func _getPlatormCorp() {
//	ctx := context.Background()
//	resp, err := eel.NewResource(ctx).LoginAsManager().Get("gskep", "corp.platform_corps", eel.Map{})
//
//	if err != nil {
//		eel.Logger.Error(err)
//		panic(err)
//	}
//
//	respData := resp.Data()
//	corp := respData.Get("corps").MustArray()[0].(map[string]interface{})
//	_PLATFORM_CORP = corp
//	corpId, _ := corp["id"].(json.Number).Int64()
//	_PLATFORM_CORP_ID = int(corpId)
//}

func init() {
	//初始化时，向gskep请求platform信息
	//_getPlatormCorp()
}
