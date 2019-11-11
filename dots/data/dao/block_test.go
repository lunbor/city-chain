package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
)

func TestBlock_Get(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123456@/citychain?charset=utf8mb4&parseTime=True")
	db.LogMode(true)
	db.SingularTable(true)

	dao := &Block{}
	m := dao.Get(db)
	if err != nil {
		t.Error(err)
	} else if m.Id != BlockId {
		t.Error("记录只有一条")
	}
}
