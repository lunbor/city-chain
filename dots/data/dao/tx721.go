package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/scryinfo/citychain/dots/data/model"
	"github.com/scryinfo/scryg/sutils/uuid"
)

type Tx721 struct {
}

func (c *Tx721) Save(db *gorm.DB, m *model.Tx721) error {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.UpdatingInit()
	return db.Save(m).Error
}

func (c *Tx721) Create(db *gorm.DB, m *model.Tx721) error {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}

	m.CreatingInit()
	return db.Save(m).Error
}
