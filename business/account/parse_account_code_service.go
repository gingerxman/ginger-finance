package account

import (
	"context"
	"github.com/gingerxman/eel"
	"strconv"
	"strings"
)

// 目前账户的Code有以下几种形式
// 1. user_cash_xx
// 2. cash
// 3. cash.xx
// 4. cash.coupon.xx
type ParseAccountCodeService struct{
	eel.ServiceBase
}

// ParseUserIdFromCode 从code中解析出user_id
func (this *ParseAccountCodeService) ParseUserIdFromCode(code string) int{
	sps := strings.Split(code, "_")
	userId, _ := strconv.Atoi(sps[len(sps)-1])
	return userId
}

// ParseImoneyCodeFromCode 从code中解析出imoney_code
func (this *ParseAccountCodeService) ParseImoneyCodeFromCode(code string) string{
	items := strings.Split(code, "_")
	subItems := items[:len(items)-1]
	if len(subItems) == 1 {
		return subItems[0]
	} else {
		return strings.Join(subItems, "_")
	}
}

func NewParseAccountCodeService(ctx context.Context) *ParseAccountCodeService{
	instance := new(ParseAccountCodeService)
	instance.Ctx = ctx
	return instance
}