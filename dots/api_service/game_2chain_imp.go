package api_service

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/scryinfo/citychain/api/server"
	"github.com/scryinfo/citychain/dots/data"
	"github.com/scryinfo/citychain/dots/data/model"
	"github.com/scryinfo/citychain/dots/erc20"
	"github.com/scryinfo/citychain/dots/erc721"
	"math/big"
)

const (
	ERR_PARAMETER_IN = "输入参数错误，可能是json不正确或没有值"
)

type Game2ChainImp struct {
	Erc20      *erc20.Erc20   `dot:""`
	Erc721     *erc721.Erc721 `dot:""`
	DataAccess *data.Access   `dot:""`
}

func (c *Game2ChainImp) GenerateErc20Tx(req *server.GenerateTransferErc20TxReq) *server.GenerateTransferErc20TxRes {
	res := &server.GenerateTransferErc20TxRes{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = ERR_PARAMETER_IN
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData

	purchase, err := c.DataAccess.Purchase.GetMaxNonce(c.DataAccess.Db, req.From)
	if err != nil {
		res.Header.Result = server.DATABASE_ERROR
		res.Header.Message = err.Error()
		return res
	}

	if purchase == nil { //第一次使用
		purchase = &model.Purchase{}
		purchase.Nonce = 0
		purchase.Status = server.Generate
		purchase.FromAddr = req.From
		purchase.ToAddr = req.To
		purchase.Contract = c.Erc20.HexContract()
		purchase.Done = false
		purchase.PurchaseId = req.Header.Id
		//这里不需要保存一次，因为最后会保存
	} else {
		if purchase.Status != server.Cancel && purchase.Status != server.Success && purchase.Status != server.Failed {
			err = errors.New("交易没有完成，不熊创建新的交易") //给数据库事务
			res.Header.Result = server.TX_NOT_FINISH
			res.Header.Message = err.Error()
			return res
		} else {
			purchase.Nonce += 1 //nonce值加1
		}
	}

	value := big.NewInt(0)
	value.SetString(req.Value, 10)
	js, err := c.Erc20.GenerateTransferTx(req.From, req.To, value, req.Header.Id, purchase.Nonce, nil, nil)
	if err != nil {
		res.Header.Result = server.GENERATE_FAIL
		res.Header.Message = err.Error()
		return res
	} else {
		res.Tx = js
		purchase.TxJson = js
		err = c.DataAccess.Purchase.Save(c.DataAccess.Db, purchase)
		if err != nil {
			res.Header.Result = server.DATABASE_ERROR
			res.Header.Message = err.Error()
			return res
		}
	}

	return res
}

func (c *Game2ChainImp) BalanceOfErc20(req *server.BalanceOfErc20Req) *server.BalanceOfErc20Res {

	res := &server.BalanceOfErc20Res{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = "输入参数错误，可能是json不正确或没有值"
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData

	address := common.HexToAddress(req.Address)
	v, err := c.Erc20.BalanceOf(&address)
	if err == nil {
		res.Value = v.String()
	} else {
		res.Header.Result = server.GETDATA_FAIL
		res.Header.Message = err.Error()
	}

	return res
}

func (c *Game2ChainImp) GenerateErc721Tx(req *server.GenerateTransferErc721TxReq) *server.GenerateTransferErc721TxRes {
	res := &server.GenerateTransferErc721TxRes{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = ERR_PARAMETER_IN
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData
	//todo
	return res
}

func (c *Game2ChainImp) BalanceOfErc721(req *server.BalanceOfErc721Req) *server.BalanceOfErc721Res {
	res := &server.BalanceOfErc721Res{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = ERR_PARAMETER_IN
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData
	//todo
	return res
}

func (c *Game2ChainImp) OwnerOfErc721(req *server.OwnerOfErc721Req) *server.OwnerOfErc721Res {
	res := &server.OwnerOfErc721Res{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = ERR_PARAMETER_IN
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData
	//todo
	return res
}

func (c *Game2ChainImp) IncreaseErc721(req *server.CreateErc721Req) *server.CreateErc721Res {
	res := &server.CreateErc721Res{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = ERR_PARAMETER_IN
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData
	//todo
	return res
}

func (c *Game2ChainImp) SendTx(req *server.SendTxReq) *server.SendTxRes {
	res := &server.SendTxRes{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = ERR_PARAMETER_IN
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData
	//todo
	return res
}

func (c *Game2ChainImp) CancelTxMust(req *server.CancelTxMustReq) *server.CancelTxMustRes {
	res := &server.CancelTxMustRes{Header: *server.NewResHeader()}
	if req == nil {
		res.Header.Result = server.PARAMETER_IN
		res.Header.Message = ERR_PARAMETER_IN
		return res
	}
	res.Header.Id = req.Header.Id
	res.Header.ExData = req.Header.ExData
	//todo
	return res
}
