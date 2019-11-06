package routers

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/handler/rest/console"
	"github.com/gingerxman/eel/handler/rest/op"
	"github.com/gingerxman/ginger-finance/rest/dev"
	"github.com/gingerxman/ginger-finance/rest/imoney"
)

func init() {
	eel.RegisterResource(&console.Console{})
	eel.RegisterResource(&op.Health{})
	
	/*
	 imoney
	 */
	eel.RegisterResource(&imoney.Imoney{})
	eel.RegisterResource(&imoney.Transfer{})
	eel.RegisterResource(&imoney.Balance{})
	
	/*
	 dev
	 */
	eel.RegisterResource(&dev.BDDReset{})
}