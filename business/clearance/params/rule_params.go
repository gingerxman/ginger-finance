package params

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	
	"github.com/gingerxman/eel"
)

type ExtraData struct{
	Role string `json:"role,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

type Validation struct{
	Type string `json:"type"`
	CorpId int `json:"corp_id,omitempty"`
	UserId int `json:"user_id,omitempty"`
}


type RuleParams struct{
	UserId interface{} `json:"user_id"`
	Name string `json:"name"`
	ExtraData *ExtraData `json:"extra_data"`
	Parent *RuleParams `json:"parent"`

	Ratio float64 `json:"ratio"`
	RelevantUsers []*RuleParams `json:"relevant_users"`
}
func NewClearanceRuleParamsFromJson(rules *simplejson.Json) []*RuleParams{

	ruleParams := make([]*RuleParams, 0)
	if rules != nil{
		rulesBytes, _ := rules.MarshalJSON()
		err := json.Unmarshal(rulesBytes, &ruleParams)
		if err != nil{
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("clearance_rule_params:decode_failed", "解析清算规则数据失败"))
		}
	}

	return ruleParams
}

func NewClearanceRuleParamsFromMap(rules []map[string]interface{}) []*RuleParams{
	ruleParams := make([]*RuleParams, 0)
	for _, rule := range rules{
		parentData := rule["parent"].(map[string]interface{})
		parent := &RuleParams{
			UserId: parentData["user_id"].(int),
			Name: parentData["name"].(string),
			ExtraData: &ExtraData{
				Role: parentData["extra_data"].(map[string]interface{})["role"].(string),
			},
		}
		ruleParam := &RuleParams{
			UserId: rule["user_id"].(int),
			Name: rule["name"].(string),
			Ratio: rule["ratio"].(float64),
			ExtraData: &ExtraData{
				Role: rule["extra_data"].(map[string]interface{})["role"].(string),
			},
			Parent: parent,
		}

		ruleParams = append(ruleParams, ruleParam)
	}
	return ruleParams
}

// FormattedRule 处理后的规则数据
// todo: 当前这个field有2种概念：
// 1、在根部(非children),表示父规则的account_id
// 2、在children中，表示当前node的account_id
// 以上需要优化，以免混淆
type FormattedRule struct {
	AccountId interface{} `json:"account_id"` // account_id可能带变量占位符，如"{related_with}"
	Ratio float64 `json:"ratio"`
	ExtraData *ExtraData `json:"extra_data,omitempty"`
	Children []*FormattedRule `json:"children,omitempty"`
}