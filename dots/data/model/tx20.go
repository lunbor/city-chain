package model

type Tx20 struct {
	Base
	TxBase
	Value       string `gorm:"not null;size:80"`
	BlockNumber string `gorm:"not null;size:80"`
}
