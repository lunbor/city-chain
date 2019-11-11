package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/scryinfo/citychain/dots/data/dao"
	"github.com/scryinfo/citychain/dots/data/model"
)

func main() {
	//&parseTime=True&loc=Local
	db, _ := gorm.Open("mysql", "root:123456@/citychain?charset=utf8mb4&parseTime=True")
	db.LogMode(true)
	db.SingularTable(true)
	db.AutoMigrate(&model.Block{}, &model.Purchase{}, &model.Tx20{}, &model.Tx721{})

	ba := &dao.Block{}
	m := ba.Get(db)
	ba.Update(db, m)
}
