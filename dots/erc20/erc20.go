package erc20

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/scryinfo/citychain/dots/data/eth"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"math/big"
	"strings"
)

const TypeId = "86dc2b53-d8bf-46fb-b1d4-384d08191cd4"

type config struct {
	Contract string `json:"contract"`
	Price    int64  `json:"price"`
	Limit    int64  `json:"limit"`
}

type Erc20 struct {
	EthConnect *eth.EthConnect `dot:""`
	abi        abi.ABI
	tranfer    abi.Method
	contract   common.Address
	cityErc20  *CityErc20

	conf  config
	price *big.Int
	limit *big.Int
}

func (c *Erc20) AfterAllInject(l dot.Line) {
	token, err := NewCityErc20(c.contract, c.EthConnect.EthClient)
	if err != nil {
		dot.Logger().Errorln("Erc20", zap.Error(err))
	} else {
		c.cityErc20 = token
	}
	c.price = big.NewInt(c.conf.Price)
	c.limit = big.NewInt(c.conf.Limit)
}

type TransferInputs struct {
	To    common.Address
	Value *big.Int
	Id    string
}

func (c *Erc20) HexContract() string {
	return c.conf.Contract
}

//get run Token input TO address
func (c *Erc20) TransferData(input []byte) (*TransferInputs, error) {
	encodedData := input[4:] //去掉方法名
	inputs := TransferInputs{}
	err := c.tranfer.Inputs.Unpack(&inputs, encodedData)
	if err != nil {
		dot.Logger().Errorln("", zap.Error(err))
		return nil, err
	} else {
		encodedData = encodedData[c.tranfer.Inputs.LengthNonIndexed()*32:]
		inputs.Id = string(encodedData)
		return &inputs, nil
	}
}

func (c *Erc20) BalanceOf(address *common.Address) (*big.Int, error) {
	if c.cityErc20 != nil {
		v, err := c.cityErc20.BalanceOf(nil, *address)
		if err != nil {
			dot.Logger().Errorln("Erc20", zap.Error(err))
		}
		return v, err
	}
	return nil, errors.New("the erc20 is null")
}

func (c *Erc20) GenerateTransferTx(from string, to string, value *big.Int, purchaseId string, nonce uint64, price *big.Int, limit *big.Int) (string, error) {
	if price == nil {
		price = c.price
	}
	if limit == nil {
		limit = c.limit
	}

	toa := common.HexToAddress(to)

	data, err := c.abi.Pack("transfer", toa, value)
	if err != nil {
		return "", err
	}
	data = append(data, ([]byte)(purchaseId)...)
	tx := types.NewTransaction(nonce, c.contract, nil, limit.Uint64(), price, data)
	bs, err := tx.MarshalJSON()

	return string(bs), err
}

func newErc20(conf interface{}) (d dot.Dot, err error) {
	var bs []byte
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &config{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	if len(dconf.Contract) < 1 {
		return nil, errors.New("the contract  is null")
	}

	erc := &Erc20{conf: *dconf}
	{
		erc.contract = common.HexToAddress(erc.conf.Contract)
		erc.abi, err = abi.JSON(strings.NewReader(CityErc20ABI))
		if err != nil {
			erc = nil
		}
		erc.tranfer = erc.abi.Methods["transfer"]
	}

	d = erc
	return d, err
}

//TypeLives
func TypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeId: TypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return newErc20(conf)
			}},
		},
	}
	lives = append(lives, eth.EthConnectTypeLives()...)
	return lives
}

//ConfigTypeLives
func ConfigTypeLives() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIdConfig: TypeId,
		ConfigInfo:   &config{},
	}
}
