package api_service

import (
	"github.com/scryinfo/citychain/api/server"
	"github.com/ybbus/jsonrpc"
	"testing"
)

func TestGame2ChainImp_BalanceOfErc20(t *testing.T) {
	req := &server.BalanceOfErc20Req{}
	res := &server.BalanceOfErc20Res{}
	//u, _ := url.Parse(c.conf.Url)
	client := jsonrpc.NewClient("http://127.0.0.1:8899")
	//server.Game2ChainName
	req.Address = "0xaa638fcA332190b63Be1605bAeFDE1df0b3b031e"
	err := client.CallFor(&res, server.Game2ChainName+".BalanceOfErc20", nil)
	if err != nil {
		t.Error(err)
	}
}
