package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
)

func TestPurchase_GetMaxNonce(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123456@/citychain?charset=utf8mb4&parseTime=True")
	db.LogMode(true)
	db.SingularTable(true)

	dao := &Purchase{}

	purchase, err := dao.GetMaxNonce(db, "no")
	if err != nil {
		t.Error(err)
	} else if purchase != nil {
		t.Error("没有找到数据时，返回空")
	}
}
