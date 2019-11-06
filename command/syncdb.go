package main

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-finance/models"
)

func main() {
	for _, model := range eel.GetRegisteredModels() {
		eel.Logger.Infof("[db] migrate table %s", model.(eel.IModel).TableName())
	}
	models.Db.AutoMigrate(eel.GetRegisteredModels()...)
}

