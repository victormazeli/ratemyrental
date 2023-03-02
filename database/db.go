package database

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Instance *gorm.DB
var dbError error

func ConnectDB(DbUsername string, DbPassword string, DbHost string, DbPort string, DbName string) *gorm.DB {
	dsn := DbUsername + ":" + DbPassword + "@tcp" + "(" + DbHost + ":" + DbPort + ")/" + DbName + "?" + "parseTime=true&loc=Local"
	sqlDB, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("Cannot open connection")
	}
	Instance, dbError = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if dbError != nil {
		//log.Fatal(dbError)
		log.Fatal("Cannot connect to DB")
	}
	log.Println("Connected to Database!")

	return Instance

}
