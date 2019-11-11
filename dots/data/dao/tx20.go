package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/scryinfo/citychain/dots/data/model"
	"github.com/scryinfo/scryg/sutils/uuid"
)

type Tx20 struct {
}

func (c *Tx20) Save(db *gorm.DB, m *model.Tx20) error {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.UpdatingInit()
	return db.Save(m).Error
}

func (c *Tx20) Create(db *gorm.DB, m *model.Tx20) error {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.CreatingInit()
	return db.Save(m).Error
}
