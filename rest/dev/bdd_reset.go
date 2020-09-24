package dev

import (
	"github.com/gingerxman/eel"
)

type BDDReset struct {
	eel.RestResource
}

func (this *BDDReset) Resource() string {
	return "dev.bdd_reset"
}

func (this *BDDReset) SkipAuthCheck() bool {
	return true
}

func (r *BDDReset) IsForDevTest() bool {
	return true
}

func (this *BDDReset) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT":  []string{},
	}
}

func (this *BDDReset) Put(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	o := eel.GetOrmFromContext(bCtx)
	
	o.Exec("delete from system_trigger_log")
	
	o.Exec("delete from clearing_record")
	o.Exec("delete from clearing_rule")
	
	o.Exec("delete from imoney_imoney")
	
	o.Exec("delete from account_transfer")
	o.Exec("delete from account_frozen_record")
	o.Exec("delete from account_account")
	o.Exec("delete from account_balance_change_log")
	
	ctx.Response.JSON(eel.Map{})
}

