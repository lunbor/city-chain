package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/scryinfo/citychain/api/server"
	"github.com/scryinfo/citychain/dots/data/model"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/scryg/sutils/uuid"
	"go.uber.org/zap"
)

type Purchase struct {
}

//如果没有找到数据，m = nil 且 err = nil
func (c *Purchase) Get(db *gorm.DB, purchaseId string) (m *model.Purchase, err error) {
	temp := &model.Purchase{}
	result := db.First(temp, "purchase_id = ?", purchaseId)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if result.Error != nil {
		temp = nil
	}
	return temp, result.Error
}

//如果没有找到数据，m = nil 且 err = nil
func (c *Purchase) GetMaxNonce(db *gorm.DB, from string) (m *model.Purchase, err error) {
	temp := &model.Purchase{}
	result := db.Where("from_addr = ? and nonce = (select max(nonce) from purchase)", from).Find(temp)
	err = result.Error
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	} else if result.Error != nil {
		temp = nil
	}
	return temp, result.Error
}

func (c *Purchase) Save(db *gorm.DB, m *model.Purchase) error {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.UpdatingInit()
	return db.Save(m).Error
}

func (c *Purchase) Create(db *gorm.DB, m *model.Purchase) error {
	if len(m.Id) < 1 {
		m.Id = uuid.GetUuid()
	}
	m.CreatingInit()
	return db.Save(m).Error
}

func (c *Purchase) Done(db *gorm.DB, m *model.Purchase) error {
	m.Done = true
	m.UpdatingInit()
	return db.Save(m).Error
}

func (c *Purchase) NotDone(db *gorm.DB) (ms []model.Purchase) {
	ms = nil
	result := db.Where("done = ? and (status = ? or status = ? )", false, server.Failed, server.Success).Find(&ms)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		dot.Logger().Errorln("", zap.Error(result.Error))
	}
	return
}
