package server

//TxResultMust
//保证通知到
type TxResultMustReq struct {
	Header ReqHeader `json:"header"`
	Status string    `json:"status"` //取消，成功，失败, see the model.Result
}
type TxResultMustRes struct {
	Header ResHeader `json:"header"`
}

//TxResultMust end

//chain to game

const Chain2GameName = "Chain2Game"

type Chain2Game interface {
	//交易结果，保证通知；如果同一id通知大于一次时，请返回执行成功
	//json method name:  Chain2Game.TxResultMust
	TxResultMust(req *TxResultMustReq) *TxResultMustRes
}
