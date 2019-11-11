package routers

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/handler/rest/console"
	"github.com/gingerxman/eel/handler/rest/op"
	"github.com/gingerxman/ginger-finance/rest/account"
	"github.com/gingerxman/ginger-finance/rest/clearance"
	"github.com/gingerxman/ginger-finance/rest/dev"
	"github.com/gingerxman/ginger-finance/rest/imoney"
)

func init() {
	eel.RegisterResource(&console.Console{})
	eel.RegisterResource(&op.Health{})
	
	/*
	  account
	 */
	eel.RegisterResource(&account.Account{})
	
	/*
	 imoney
	 */
	eel.RegisterResource(&imoney.Imoney{})
	eel.RegisterResource(&imoney.Transfer{})
	eel.RegisterResource(&imoney.Balance{})
	eel.RegisterResource(&imoney.FrozenRecord{})
	
	/*
	 clearance
	 */
	eel.RegisterResource(&clearance.OrderClearance{})
	
	/*
	 dev
	 */
	eel.RegisterResource(&dev.BDDReset{})
}