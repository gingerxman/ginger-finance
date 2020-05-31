package imoney

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/config"
	m_imoney "github.com/gingerxman/ginger-finance/models/imoney"
)

var Imoneys []*Imoney
var Code2Imoney map[string]*Imoney

type ImoneyManager struct {
	eel.ServiceBase
}

func (this *ImoneyManager) ImoneyExisted(code string) bool{
	if _, ok := Code2Imoney[code]; ok{
		return true
	}
	return false
}

func (this *ImoneyManager) GetImoneyByCode(code string) *Imoney{
	if imoney, ok := Code2Imoney[code]; ok{
		return imoney
	}
	return nil
}

func (this *ImoneyManager) Add(imoney *Imoney){
	if this.ImoneyExisted(imoney.Code){
		panic(eel.NewBusinessError("imoney_existed", "虚拟资产已存在"))
	}
	o := eel.GetOrmFromContext(this.Ctx)
	result := o.Create(&m_imoney.IMoney{
		Code: imoney.Code,
		DisplayName: imoney.DisplayName,
		ExchangeRate: imoney.ExchangeRate,
		IsPayable: imoney.IsPayable,
		IsDebtable: imoney.IsDebtable,
	})
	if err := result.Error; err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("imoney:save_failed", "存储imoney失败"))
	}
	Imoneys = append(Imoneys, imoney)
	Code2Imoney[imoney.Code] = imoney
}

func (this *ImoneyManager) Remove(imoney *Imoney){
	if this.ImoneyExisted(imoney.Code){
		return
	}
	o := eel.GetOrmFromContext(this.Ctx)
	result := o.Where("code", imoney.Code).Delete(m_imoney.IMoney{})
	if err := result.Error; err != nil{
		panic(eel.NewBusinessError("imoney:delete_failed", "删除虚拟资产失败"))
	}
	for _, im := range Imoneys{
		target := Imoneys[:0] // 共用原始数组，达到原地修改的目的
		if im.Code != imoney.Code{
			target = append(target, im)
		}
	}
	delete(Code2Imoney, imoney.Code)
}

func NewImoneyManager(ctx context.Context) *ImoneyManager{
	inst := new(ImoneyManager)
	inst.Ctx = ctx
	return inst
}

// init 从数据库中加载imoney
func init() {
	Imoneys = make([]*Imoney, 0)
	Code2Imoney = make(map[string]*Imoney)

	o := config.Runtime.DB
	var dbModels []*m_imoney.IMoney
	result := o.Model(m_imoney.IMoney{}).Find(&dbModels)
	if err := result.Error; err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("imoney:load_failed", "加载虚拟资产失败"))
	}
	for _, dbModel := range dbModels{
		imoney := NewImoneyFromModel(dbModel)
		Imoneys = append(Imoneys, imoney)
		Code2Imoney[dbModel.Code] = imoney
	}
}