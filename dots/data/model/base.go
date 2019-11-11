package model

import (
	"fmt"
	"github.com/scryinfo/citychain/api/server"
	"time"
)

type Base struct {
	Id        string     `gorm:"primary_key;size:36"`
	Version   int64      `gorm:"default:1;not null"`
	CreatedAt int64      `gorm:"not null"` //; DEFAULT:CURRENT_TIMESTAMP
	UpdatedAt int64      `gorm:"not null"` //; DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	DeletedAt *time.Time `sql:"index"`     //这是gorm的特殊字段， 在select /update时会主动增加条件 deleted_at IS NULL
}

func (c *Base) CreateAtTime() time.Time {
	return time.Unix(c.CreatedAt, 0)
}
func (c *Base) UpdatedAtTime() time.Time {
	return time.Unix(c.UpdatedAt, 0)
}
func (c *Base) CreatingInit() {
	c.CreatedAt = time.Now().Unix()
	c.UpdatedAt = c.CreatedAt
}

func (c *Base) UpdatingInit() {
	c.UpdatedAt = time.Now().Unix()
}

type TxBase struct {
	PurchaseId string        `gorm:"not null;size:36;unique_index:uq_purchase"` //保证purchase id是唯一的
	TxId       string        `gorm:"not null;size:66"`
	FromAddr   string        `gorm:"not null;size:42;unique_index:uq_nonce"` //保证 from nonce是唯一的
	ToAddr     string        `gorm:"not null;size:42"`                       //不是合约地址
	Contract   string        `gorm:"not null;size:42"`                       //调用的合约地址
	Nonce      uint64        `gorm:"not null;unique_index:uq_nonce"`
	Status     server.Status `gorm:"not null;size:16;default:'Generate'"` //保证 from nonce是唯一的
}

func (c *TxBase) String() string {
	return fmt.Sprintf("purchase id: %s, tx id: %s, from: %s, to: %s, status: %s", c.PurchaseId, c.TxId, c.FromAddr, c.ToAddr, c.Status)
}
