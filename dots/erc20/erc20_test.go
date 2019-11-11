package erc20

import (
	"encoding/json"
	"github.com/scryinfo/citychain/dots/data/eth"
	"math/big"
	"testing"
)

func TestErc20_GenerateTransferTx(t *testing.T) {

	var erc20 *Erc20
	{
		conf := &config{
			Contract: "0xaa638fcA332190b63Be1605bAeFDE1df0b3b031e",
			Price:    200000,
			Limit:    600000,
		}
		b, _ := json.Marshal(conf)
		d, _ := newErc20(b)
		erc20, _ = d.(*Erc20)
		erc20.EthConnect = &eth.EthConnect{}
		erc20.AfterAllInject(nil)
	}

	js, _ := erc20.GenerateTransferTx("", "0xaa638fcA332190b63Be1605bAeFDE1df0b3b0316", big.NewInt(16), "test", 1, nil, nil)
	if len(js) < 1 {
		t.Error("")
	}
}
