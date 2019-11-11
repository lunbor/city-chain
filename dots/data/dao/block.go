package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/scryinfo/citychain/dots/data/model"
	"math/big"
)

const BlockId = "1"

const (
	NotGet   = "-1"
	GetError = "-2"
)

var (
	NotGetBigInt   = big.NewInt(-1)
	GetErrorBigInt = big.NewInt(-2)
)

type Block struct {
}

func (c *Block) Update(db *gorm.DB, m *model.Block) error {
	if m.Id != BlockId {
		m.Id = BlockId
	}
	m.UpdatingInit()
	result := db.Save(m)
	return result.Error
}

//如果没有找到，就返回一条默认的
func (c *Block) Get(db *gorm.DB) *model.Block {
	var temp model.Block
	result := db.Where("id = ?", BlockId).First(&temp)
	if result.Error != nil || temp.Id != BlockId {
		temp.DoneBlock = NotGet
		temp.Id = BlockId
		temp.Version = 0
		temp.CreatingInit()
	}
	return &temp
}
