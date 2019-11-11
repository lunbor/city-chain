package server

//GenerateTransferErc20Tx
type GenerateTransferErc20TxReq struct {
	Header ReqHeader `json:"header"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Value  string    `json:"value"` //数量是一个大整数类型，“1,000,000,000,000,000,000” == 1个 erc20
}
type GenerateTransferErc20TxRes struct {
	Header ResHeader `json:"header"`
	Tx     string    `json:"tx"` //没有签名的交易
}

//GenerateTransferErc20Tx end

//BalanceOfErc20
type BalanceOfErc20Req struct {
	Header  ReqHeader `json:"header"`
	Address string    `json:"address"`
}
type BalanceOfErc20Res struct {
	Header ResHeader `json:"header"`
	Value  string    `json:"value"` //数量是一个大整数类型，“1,000,000,000,000,000,000” == 1个 erc20
}

//BalanceOfErc20 end

//GenerateTransferErc721Tx
type GenerateTransferErc721TxReq struct {
	Header  ReqHeader `json:"header"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	TokenId uint64    `json:"tokenId"`
}
type GenerateTransferErc721TxRes struct {
	Header ResHeader `json:"header"`
	Tx     string
}

//GenerateTransferErc721Tx end

//BalanceOfErc721
type BalanceOfErc721Req struct {
	Header  ReqHeader `json:"header"`
	Address string    `json:"address"`
}
type BalanceOfErc721Res struct {
	Header ResHeader `json:"header"`
	Value  uint64    `json:"value"` //token id的数量
}

//BalanceOfErc721 end

//OwnerOfErc721
type OwnerOfErc721Req struct {
	Header  ReqHeader `json:"header"`
	TokenId uint64    `json:"tokenId"`
}
type OwnerOfErc721Res struct {
	Header  ResHeader `json:"header"`
	Address string    `json:"address"`
}

//OwnerOfErc721 end

//CreateErc721
type CreateErc721Req struct {
	Header  ReqHeader `json:"header"`
	TokenId uint64    `json:"tokenId"`
}
type CreateErc721Res struct {
	Header ResHeader `json:"header"`
	Tx     string    `json:"tx"`
}

//CreateErc721 end

//DestroyErc721Req
type DestroyErc721Req struct {
	Header  ReqHeader `json:"header"`
	TokenId uint64    `json:"tokenId"`
}
type DestroyErc721Res struct {
	Header ResHeader `json:"header"`
	Tx     string    `json:"tx"`
}

//DestroyErc721Req end

//SendTx
type SendTxReq struct {
	Header   ReqHeader `json:"header"`
	SingedTx string    `json:"singedTx"` //要转为json，使用字符串更为合理
}
type SendTxRes struct {
	Header ResHeader `json:"header"`
	TxId   string    `json:"txId"`
}

//SendTx end

//CancelTxMust
//保证通知到
type CancelTxMustReq struct {
	Header ReqHeader `json:"header"`
	Note   string    `json:"note"` //取消说明，只是做记录，字符不超过16字符
}
type CancelTxMustRes struct {
	Header ResHeader `json:"header"`
}

//CancelTxMust end

//game to chain
const Game2ChainName = "Game2Chain"

type Game2Chain interface {
	//生成erc20的转帐交易，需要订单id
	//json method name:  Game2Chain.GenerateTransferErc20Tx
	GenerateTransferErc20Tx(req *GenerateTransferErc20TxReq) *GenerateTransferErc20TxRes
	//查询erc20的余额
	//json method name:  Game2Chain.BalanceOfErc20
	BalanceOfErc20(req *BalanceOfErc20Req) *BalanceOfErc20Res

	////创建erc20
	////json method name:  Game2Chain.CreateErc20
	//CreateErc20(req *CreateErc20Req) *CreateErc20Res
	////销毁erc20
	////json method name:  Game2Chain.DestroyErc20
	//DestroyErc20(req *DestroyErc20Req) *DestroyErc20Res

	//生成erc721的转帐交易，需要订单id
	//json method name:  Game2Chain.GenerateTransferErc721Tx
	GenerateTransferErc721Tx(req *GenerateTransferErc721TxReq) *GenerateTransferErc721TxRes
	//查询erc721的余额
	//json method name:  Game2Chain.BalanceOfErc721
	BalanceOfErc721(req *BalanceOfErc721Req) *BalanceOfErc721Res
	//查询erc721 token id对应的拥有者地址
	//json method name:  Game2Chain.OwnerOfErc721
	OwnerOfErc721(req *OwnerOfErc721Req) *OwnerOfErc721Res
	//创建erc721，需要订单id
	//json method name:  Game2Chain.CreateErc721
	CreateErc721(req *CreateErc721Req) *CreateErc721Res
	//销毁erc721，需要订单id
	//json method name:  Game2Chain.CreateErc721
	DestroyErc721(req *DestroyErc721Req) *DestroyErc721Res

	//发送交易
	//json method name:  Game2Chain.SendTx
	SendTx(req *SendTxReq) *SendTxRes
	//取消交易，保证通知到。 生成交易之后，发送交易之前，可以取消
	//json method name:  Game2Chain.CancelTxMust
	CancelTxMust(req *CancelTxMustReq) *CancelTxMustRes
}
