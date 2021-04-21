package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	//"github.com/jinzhu/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hita/config"
	"log"
)

//因为我们需要在其他地方使用SqlDB这个变量，所以需要大写代表public
var DB *gorm.DB

//初始化方法
func InitDB() error {
	var err error
	dbDriver := config.MysqlUsername + ":" + config.MysqlPassword + "@tcp(" + config.MysqlIp + ":" + config.MysqlPort + ")/" + config.MysqlDbname + "?charset=utf8&parseTime=true"
	DB, err = gorm.Open(mysql.Open(dbDriver))
	//DB, err = gorm.Open("mysql", dbDriver)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
