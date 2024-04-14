package model

import (
	"fmt"
	"go-file/common"
	"go-file/common/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func createAdminAccount() {
	var user User
	DB.Where(User{Role: common.RoleAdminUser}).Attrs(User{
		Username:    "admin",
		Password:    "123456",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		DisplayName: "Administrator",
	}).FirstOrCreate(&user)
}

func CountTable(tableName string) (num int) {
	DB.Table(tableName).Count(&num)
	return
}

func InitDB(conf *config.Config) (db *gorm.DB, err error) {
	fmt.Println("sqldns", conf.SqlDNS)
	if conf.SqlDNS != "" {
		// Use MySQL
		db, err = gorm.Open("mysql", conf.SqlDNS)
	} else {
		// Use SQLite
		db, err = gorm.Open("sqlite3", common.SQLitePath)
	}
	if err == nil {
		DB = db
		db.AutoMigrate(&File{})
		db.AutoMigrate(&Image{})
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Option{})
		createAdminAccount()
		return DB, err
	} else {
		common.FatalLog("failed to connect to database: " + err.Error())
	}
	return nil, err
}
