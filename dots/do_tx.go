package dots

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/scryinfo/citychain/api/server"
	"github.com/scryinfo/citychain/dots/api_service"
	"github.com/scryinfo/citychain/dots/data"
	"github.com/scryinfo/citychain/dots/data/dao"
	"github.com/scryinfo/citychain/dots/data/model"
	"github.com/scryinfo/citychain/dots/erc20"
	"github.com/scryinfo/citychain/dots/erc721"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"math/big"
)

const DoTxTypeId = "3cca7d6c-779f-484a-b2fa-c56f8ef9eed2"

type doTxConfig struct {
}

type DoTx struct {
	Scan       *Scan                         `dot:""`
	DataAccess *data.Access                  `dot:""`
	Chain2Game *api_service.Chain2GameClient `dot:""`
	Erc20      *erc20.Erc20                  `dot:""`
	Erc721     *erc721.Erc721                `dot:""`

	conf doTxConfig
}

func (c *DoTx) StartBlockNumber() *big.Int {
	if c.DataAccess != nil {
		b := c.DataAccess.Block.Get(c.DataAccess.Db)
		if b == nil {
			return dao.GetErrorBigInt
		} else {
			d := big.NewInt(0)
			d.SetString(b.DoneBlock, 10)
			return d
		}
	} else {
		return dao.GetErrorBigInt
	}
}

func (c *DoTx) Tx(bl *types.Block, tx *types.Transaction, receipt *types.Receipt) bool {

	contract := tx.To().Hex()
	switch contract {
	case c.Erc20.HexContract():
		return c.TxErc20(bl, tx, receipt)
	case c.Erc721.HexContract():
		return c.TxErc721(bl, tx, receipt)
	default:
		//是不关注的交易，直接返回成功
	}

	return true
}

func (c *DoTx) TxErc20(bl *types.Block, tx *types.Transaction, receipt *types.Receipt) bool {
	logger := dot.Logger()
	var err error
	for {
		contract := tx.To().Hex()
		tx20 := model.Tx20{}
		tx20.Contract = contract
		tx20.BlockNumber = bl.Number().String()
		tx20.TxId = tx.Hash().Hex()
		if receipt.Status == 1 {
			tx20.Status = server.Success
		} else {
			tx20.Status = server.Failed
		}
		{
			var sender common.Address
			sender, err = c.Scan.EthConnect.Signer.Sender(tx)
			if err != nil {
				break
			}
			tx20.FromAddr = sender.Hex()
		}
		{
			var inputs *erc20.TransferInputs
			inputs, err = c.Erc20.TransferData(tx.Data())
			if err != nil {
				break
			}
			tx20.Value = inputs.Value.String()
			tx20.ToAddr = inputs.To.Hex()
			tx20.PurchaseId = inputs.Id
		}

		tx20.Nonce = tx.Nonce()

		err = c.DataAccess.Tx20.Save(c.DataAccess.Db, &tx20)
		if err != nil {
			break
		}
		var purchase *model.Purchase
		purchase, err = c.DataAccess.Purchase.Get(c.DataAccess.Db, tx20.PurchaseId)
		if err != nil {
			//todo 已经有了交易，但是还没有产生订单， 这是一种异常情况；可以考虑是否直接把订单保存下来
			logger.Errorln("DoTx", zap.String("", "erc20 链上的交易没有对应的订单： "+tx20.String()))
			break
		}
		purchase.Status = tx20.Status
		err = c.DataAccess.Purchase.Save(c.DataAccess.Db, purchase)
		if err != nil {
			break
		}

		c.Chain2Game.Append(purchase)

		break
	}

	if err != nil {
		logger.Errorln("DoTx", zap.Error(err))
		return false
	}

	return true
}
func (c *DoTx) TxErc721(bl *types.Block, tx *types.Transaction, receipt *types.Receipt) bool {
	logger := dot.Logger()
	var err error
	for {
		contract := tx.To().Hex()
		tx721 := model.Tx721{}
		tx721.Contract = contract
		tx721.BlockNumber = bl.Number().String()
		tx721.TxId = tx.Hash().Hex()
		if receipt.Status == 1 {
			tx721.Status = server.Success
		} else {
			tx721.Status = server.Failed
		}
		{
			var sender common.Address
			sender, err = c.Scan.EthConnect.Signer.Sender(tx)
			if err != nil {
				break
			}
			tx721.FromAddr = sender.Hex()
		}
		{
			var inputs *erc721.TransferInputs
			inputs, err = c.Erc721.TransferData(tx.Data())
			if err != nil {
				break
			}
			tx721.TokenId = inputs.TokenId
			tx721.ToAddr = inputs.To.Hex()
			tx721.PurchaseId = inputs.Id
		}
		tx721.Nonce = tx.Nonce()

		err = c.DataAccess.Tx721.Save(c.DataAccess.Db, &tx721)
		if err != nil {
			break
		}

		var purchase *model.Purchase
		purchase, err = c.DataAccess.Purchase.Get(c.DataAccess.Db, tx721.PurchaseId)
		if err != nil {
			//todo 已经有了交易，但是还没有产生订单， 这是一种异常情况；可以考虑是否直接把订单保存下来
			logger.Errorln("DoTx", zap.String("", "erc721 链上的交易没有对应的订单： "+tx721.String()))
			break
		}
		purchase.Status = tx721.Status
		err = c.DataAccess.Purchase.Save(c.DataAccess.Db, purchase)
		if err != nil {
			break
		}
		c.Chain2Game.Append(purchase)

		break
	}

	if err != nil {
		logger.Errorln("DoTx", zap.Error(err))
		return false
	}
	return true
}
func (c *DoTx) Block(bl *types.Block) bool {
	return true
}

func (c *DoTx) DoneBlock(bl *types.Block) bool {
	dot.Logger().Debugln("DoTx", zap.String("", fmt.Sprint(bl.Header().Number)))
	b := c.DataAccess.Block.Get(c.DataAccess.Db)
	b.DoneBlock = bl.Header().Number.String()
	err := c.DataAccess.Block.Update(c.DataAccess.Db, b)
	return err == nil
}

func newDoTx(conf interface{}) (d dot.Dot, err error) {
	var bs []byte
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &doTxConfig{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	doTx := &DoTx{conf: *dconf}
	d = doTx
	return d, err
}

//DoTxTypeLives
func DoTxTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeId: DoTxTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return newDoTx(conf)
			}},
		},
	}
	lives = append(lives, ScanTypeLives()...)
	lives = append(lives, api_service.Chain2GameClientTypeLives()...)
	lives = append(lives, erc20.TypeLives()...)
	lives = append(lives, erc721.TypeLives()...)
	return lives
}

//DoTxConfigTypeLives
func DoTxConfigTypeLives() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIdConfig: DoTxTypeId,
		ConfigInfo:   &doTxConfig{},
	}
}
