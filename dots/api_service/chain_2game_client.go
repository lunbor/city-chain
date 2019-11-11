package api_service

import (
	"github.com/scryinfo/citychain/api/server"
	"github.com/scryinfo/citychain/dots/data"
	"github.com/scryinfo/citychain/dots/data/model"
	"github.com/scryinfo/dot/dot"
	"github.com/ybbus/jsonrpc"
	"go.uber.org/zap"
)

const Chain2GameClientTypeId = "65533eac-664b-45e7-996c-54f1b6a94a9e"

type chain2GameClientConfig struct {
	Url           string `json:"url"`
	IntervalRetry int64  `json:"intervalRetry"` //重试的间隔
}

type Chain2GameClient struct {
	conf       chain2GameClientConfig
	DataAccess *data.Access `dot:""`
	pps        MustPurchases
}

func (c *Chain2GameClient) TxResultMust(req *server.TxResultMustReq) (res *server.TxResultMustRes) {
	res = &server.TxResultMustRes{Header: *server.NewResHeader()}
	//u, _ := url.Parse(c.conf.Url)
	client := jsonrpc.NewClient(c.conf.Url)
	err := client.CallFor(&res, server.Chain2GameName+"TxResultMust", req)
	if err != nil {
		res.Header.Result = server.HTTP_ERROR
		res.Header.Message = err.Error()
	}
	return res
}

//在调用这个方法前，要保证数据库
func (c *Chain2GameClient) Append(pp *model.Purchase) {
	c.pps.Append(pp)
}

func (c *Chain2GameClient) AfterAllStart(l dot.Line) {
	not := c.DataAccess.Purchase.NotDone(c.DataAccess.Db)
	if len(not) > 0 {
		c.pps.Init(not, c.conf.IntervalRetry, func(purchase *model.Purchase) bool {
			req := server.TxResultMustReq{
				Header: *server.NewReqHeader(),
			}
			if purchase.Status == server.Failed || purchase.Status == server.Success {
				req.Status = string(purchase.Status)
			} else {
				dot.Logger().Errorln("", zap.String("", "通知game时，状态只取是 成功或失败"))
				//todo exit progress ?
				return false
			}

			res := c.TxResultMust(&req)
			if res.Header.Result == server.SUCCESS_INT { //通知到了，写入数据库
				err := c.DataAccess.Purchase.Done(c.DataAccess.Db, purchase)
				if err != nil {
					dot.Logger().Errorln("", zap.Error(err))
					return false
				}
			}
			return res.Header.Result == server.SUCCESS_INT
		})
	}
}

func (c *Chain2GameClient) Stop(ignore bool) error {
	var err error = nil
	c.pps.Stop()
	return err
}

func newDoTx(conf interface{}) (d dot.Dot, err error) {
	var bs []byte
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &chain2GameClientConfig{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	d = &Chain2GameClient{conf: *dconf}

	return d, err
}

//Chain2GameClientTypeLives
func Chain2GameClientTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeId: Chain2GameClientTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return newDoTx(conf)
			}},
		},
	}
	lives = append(lives, data.TypeLives()...)
	return lives
}

//Chain2GameClientConfigTypeLives
func Chain2GameClientConfigTypeLives() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIdConfig: Chain2GameClientTypeId,
		ConfigInfo:   &chain2GameClientConfig{},
	}
}
