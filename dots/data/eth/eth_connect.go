package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/scryinfo/dot/dot"
	"math/big"
)

const (
	EthConnectTypeId = "0c73cc80-ed3c-4dc9-8d4f-bcaa4dad29cd"
)

type ethConnectConfig struct {
	UrlEth  string `json:"urlEth"` //eth连接参数
	ChainId int64  `json:"chainId"`
}

type EthConnect struct {
	conf      ethConnectConfig
	Ctx       context.Context
	cancelFun context.CancelFunc
	chainId   big.Int

	EthClient *ethclient.Client //todo 确认是否会自动重连
	Signer    types.Signer
}

func (c *EthConnect) Stop(ignore bool) error {
	if c.cancelFun != nil {
		c.cancelFun()
	}
	if c.EthClient != nil {
		c.EthClient.Close()
	}

	return nil
}

func (c *EthConnect) Create(l dot.Line) error {
	var err error
	c.Ctx, c.cancelFun = context.WithCancel(context.Background())
	c.EthClient, err = ethclient.DialContext(c.Ctx, c.conf.UrlEth)
	if err == nil {
		c.chainId = *big.NewInt(c.conf.ChainId)
		c.Signer = types.NewEIP155Signer(&c.chainId)
	}
	//todo 确认重连，及第一次连接失败时的返回值

	return err
}

func newEthConnect(conf interface{}) (d dot.Dot, err error) {
	var bs []byte
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &ethConnectConfig{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	d = &EthConnect{conf: *dconf}

	return d, err
}

//EthConnectTypeLives
func EthConnectTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeId: EthConnectTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return newEthConnect(conf)
			}},
		},
	}
	return lives
}

//EthConnectConfigTypeLives
func EthConnectConfigTypeLives() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIdConfig: EthConnectTypeId,
		ConfigInfo:   &ethConnectConfig{},
	}
}
