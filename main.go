package main

import (
	"fmt"

	"github.com/Benyam-S/onepay/dbfiles"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open("mysql", "root:0911@tcp(127.0.0.1:3306)/onepay?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to the database: mysql @GORM")

	dbfiles.Init(db)
}
