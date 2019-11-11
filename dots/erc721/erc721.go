package erc721

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/scryinfo/citychain/dots/data/eth"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/zap"
	"strings"
)

const TypeId = "a91ae735-288d-4afd-b549-b78c1d90a91b"

type config struct {
	Contract string `json:"contract"`
}

type Erc721 struct {
	EthConnect *eth.EthConnect `dot:""`

	abi      abi.ABI
	tranfer  abi.Method
	contract common.Address

	conf config
}

type TransferInputs struct {
	To      common.Address
	TokenId int64
	Id      string
}

func (c *Erc721) HexContract() string {
	return c.conf.Contract
}

//todo
func (c *Erc721) TransferData(input []byte) (*TransferInputs, error) {
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

//todo
func newErc721(conf interface{}) (d dot.Dot, err error) {
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

	erc := &Erc721{conf: *dconf}
	{
		erc.contract = common.HexToAddress(erc.conf.Contract)
		erc.abi, err = abi.JSON(strings.NewReader(ABI))
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
				return newErc721(conf)
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
