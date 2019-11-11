package model

type Tx721 struct {
	Base
	TxBase
	TokenId     int64  `gorm:"not null;unique_index:uq_token_id"` //保证token id是唯一的
	BlockNumber string `gorm:"not null;size:80"`
}
