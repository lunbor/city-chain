package api_service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
	"github.com/scryinfo/citychain/api/server"
	"github.com/scryinfo/citychain/dots/data"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/utils"
	"go.uber.org/zap"
)

const (
	TypeId = "6c5acd0f-6135-4a9c-a88b-ff7734c56a12"
)

type config struct {
	AddrGame2Chain string          `json:"addrGame2Chain"`
	Tls            utils.TlsConfig `json:"tls"`
}

type ApiService struct {
	conf       config
	DataAccess *data.Access `dot:""`

	game2Chain server.Game2Chain
	httpServer *http.Server
}

func (c *ApiService) AfterAllInject(l dot.Line) {
	l.ToInjecter().Inject(c.game2Chain)
}

func (c *ApiService) AfterAllStart(l dot.Line) {
	c.startRpc()
}

func (c *ApiService) Stop(ignore bool) error {
	var err error = nil
	if c.httpServer != nil {
		err = c.httpServer.Shutdown(context.Background())
	}
	//todo set nil ?
	return err
}

func (c *ApiService) startRpc() (err error) {
	logger := dot.Logger()

	{
		handler := jsonrpc.NewServer(c.makeJsonrpc(server.Game2ChainName + "."))
		c.httpServer = &http.Server{Handler: handler}
	}

	conn, err := net.Listen("tcp", c.conf.AddrGame2Chain)
	if err != nil {
		return err
	}

	var serverRun func() error

	switch {
	case len(c.conf.Tls.CaPem) > 0 && len(c.conf.Tls.Key) > 0 && len(c.conf.Tls.Pem) > 0: //both tls
		caPem := utils.GetFullPathFile(c.conf.Tls.CaPem)
		if len(caPem) < 1 {
			logger.Errorln("ApiService", zap.Error(errors.New("the caPem is not empty, and can not find the file: "+c.conf.Tls.CaPem)))
			return
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			logger.Errorln("ApiService", zap.Error(errors.New("the Key is not empty, and can not find the file: "+c.conf.Tls.Key)))
			return
		}

		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			logger.Errorln("ApiService", zap.Error(errors.New("the Pem is not empty, and can not find the file: "+c.conf.Tls.Pem)))
			return
		}

		pool := x509.NewCertPool()
		{
			caCrt, err1 := ioutil.ReadFile(caPem)
			if err1 != nil {
				logger.Errorln("ApiService", zap.Error(errors.WithStack(err1)))
				return
			}
			if !pool.AppendCertsFromPEM(caCrt) {
				logger.Errorln("ApiService", zap.Error(errors.New("credentials: failed to append certificates")))
				return
			}
		}
		cert, err1 := tls.LoadX509KeyPair(pem, key)
		if err1 != nil {
			logger.Errorln("ApiService", zap.Error(errors.WithStack(err1)))
			return
		}

		c.httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		logger.Infoln("ApiService", zap.String("", "ApiService server(with ca) will start: "+c.conf.AddrGame2Chain))
		serverRun = func() error {
			return c.httpServer.ServeTLS(conn, "", "")
		}
	case len(c.conf.Tls.Key) > 0 && len(c.conf.Tls.Pem) > 0: //just server
		pem := utils.GetFullPathFile(c.conf.Tls.Pem)
		if len(pem) < 1 {
			logger.Errorln("ApiService", zap.Error(errors.New("the pem is not empty, and can not find the file: "+c.conf.Tls.Pem)))
			return
		}
		key := utils.GetFullPathFile(c.conf.Tls.Key)
		if len(key) < 1 {
			logger.Errorln("ApiService", zap.Error(errors.New("the key is not empty, and can not find the file: "+c.conf.Tls.Key)))
			return
		}

		c.httpServer.TLSConfig = &tls.Config{
			ClientAuth: tls.NoClientCert,
		}

		logger.Infoln("ApiService", zap.String("", "ApiService server(no ca) will start: "+c.conf.AddrGame2Chain))
		serverRun = func() error {
			return c.httpServer.ServeTLS(conn, pem, key)
		}
	default: //no tls
		logger.Infoln("ApiService", zap.String("", "ApiService server(no https) will start: "+c.conf.AddrGame2Chain))
		serverRun = func() error {
			return c.httpServer.Serve(conn)
		}
	}

	go func() {
		e := serverRun()
		if e != nil {
			logger.Infoln("ApiService", zap.Error(e))
		}
		conn.Close()
	}()

	return nil
}

//json rpc的解码函数
func nopDecoder(ctx context.Context, j json.RawMessage) (interface{}, error) {
	return j, nil
}

//json rpc的编码函数
func nopEncoder(ctx context.Context, req interface{}) (json.RawMessage, error) {
	bs, _ := json.Marshal(req)
	return bs, nil
}

//todo 有时间重新考虑使用 eth的jsonrpc
func (c *ApiService) makeJsonrpc(preName string) jsonrpc.EndpointCodecMap {
	ecm := jsonrpc.EndpointCodecMap{}
	receiver := reflect.ValueOf(c.game2Chain)
	typ := receiver.Type()
	for i := 0; i < receiver.NumMethod(); i++ {
		tmethod := typ.Method(i)
		vmethod := receiver.Method(i)
		if tmethod.PkgPath != "" {
			continue // tmethod not exported
		}
		fn := tmethod.Type
		if fn.NumOut() == 1 && fn.NumIn() == 2 { //只处理单个参数
			endpointCodec := jsonrpc.EndpointCodec{
				Decode: nopDecoder,
				Encode: nopEncoder,
			}
			endpointCodec.Endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error) {
				defer func() {
					if e2 := recover(); e2 != nil {
						response = nil
						err = nil //todo 不太确定 err的具体工作过程， 这里暂定
						dot.Logger().Debugln("ApiService", zap.Any("", e2))
					}
				}()
				inType := fn.In(1)
				inValue := reflect.New(inType)
				bs, ok := request.(json.RawMessage)
				if ok {
					err := json.Unmarshal(bs, inValue.Interface())
					if err != nil {
						inValue = reflect.New(inType)
					}
				} else {
					inValue = reflect.New(inType)
				}
				rev := vmethod.Call([]reflect.Value{inValue.Elem()})

				response = rev[0].Interface()
				err = nil
				return
			}
			if len(preName) > 0 {
				ecm[preName+tmethod.Name] = endpointCodec
			} else {
				ecm[tmethod.Name] = endpointCodec
			}
		}
	}
	return ecm
}

func newApiService(conf interface{}) (d dot.Dot, err error) {
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

	d = &ApiService{conf: *dconf, game2Chain: &Game2ChainImp{}}

	return d, err
}

func TypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeId: TypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return newApiService(conf)
			}},
		},
	}
	lives = append(lives, data.TypeLives()...)
	return lives
}

func ConfigTypeLives() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIdConfig: TypeId,
		ConfigInfo:   &config{},
	}
}
