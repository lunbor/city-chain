package model

type Block struct {
	Base
	DoneBlock string `gorm:"not null;default:'-1';size:80"` //完成的区块号, 如果为-1，表示一个都没有完成
}
