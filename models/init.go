package models

import (
	"fmt"
	
	"github.com/gingerxman/gorm"
	_ "github.com/gingerxman/gorm/dialects/mysql"

	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/config"
	_ "github.com/gingerxman/ginger-finance/models/account"
	_ "github.com/gingerxman/ginger-finance/models/imoney"
	_ "github.com/gingerxman/ginger-finance/models/clearance"
)

var Db *gorm.DB

func init() {
	host := config.ServiceConfig.String("db::DB_HOST")
	port := config.ServiceConfig.String("db::DB_PORT")
	db := config.ServiceConfig.String("db::DB_NAME")
	user := config.ServiceConfig.String("db::DB_USER")
	password := config.ServiceConfig.String("db::DB_PASSWORD")
	charset := config.ServiceConfig.String("db::DB_CHARSET")
	mysqlURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?interpolateParams=true&charset=%s&parseTime=True&loc=Asia%%2FShanghai", user, password, host, port, db, charset)
	
	var err error
	Db, err = gorm.Open("mysql", mysqlURL)

	if err != nil {
		eel.Logger.Errorw("[db] connect to mysql fail!!", "error", err.Error())
	} else {
		eel.Logger.Infof("[db] connect to mysql %s success!", mysqlURL)
	}

	Db.LogMode(true)
	eel.Runtime.DB = Db
}
