package model

import "fmt"

type Purchase struct {
	Base
	TxBase
	TxJson string `gorm:"not null;size:10240"`    //生成的交易的json表示， 还没有签名， todo 确定最大长度
	Done   bool   `gorm:"not null;default:false"` //是否通知到 game server
}

func (c *Purchase) String() string {
	return fmt.Sprintf("purchase id: %s, tx id: %s, from: %s, to: %s, status: %s", c.PurchaseId, c.TxId, c.FromAddr, c.ToAddr, c.Status)
}
