package dots

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/scryinfo/citychain/dots/data/dao"
	"github.com/scryinfo/citychain/dots/data/eth"
	"github.com/scryinfo/dot/dot"
	"go.uber.org/atomic"
	"math/big"
	"time"
)

const ScanTypeId = "872217f5-1fd4-4696-bcfd-ab99b6b321fd"

const (
	IntervalScanCall = 10 //重试调用ScanCall出错的间隔
	IntervalBlock    = 2  //重试取得区块出错的间隔
	IntervalTx       = 2  //重试取交易出错的间隔
	FinalHight       = 4  //在遍历区块时，要小于最大区块- FinalHight
)

type ScanCall interface {
	StartBlockNumber() *big.Int                                             //上次处理过的 block number， 第一次时应该从数据库中加载， 如果返回值为-2，说明还没有正常启动,-1 没有以前的值
	Tx(bl *types.Block, tx *types.Transaction, receipt *types.Receipt) bool //新的交易
	Block(bl *types.Block) bool                                             //新的区块
	DoneBlock(bl *types.Block) bool                                         //完成一个区块
}

type Scan struct {
	ScanCall   ScanCall        `dot:""`
	EthConnect *eth.EthConnect `dot:""`

	stopped         atomic.Bool
	stopChanel      chan bool
	initBlockNumber *big.Int
}

func (c *Scan) Stop(ignore bool) error {
	c.stopped.Store(true)
	close(c.stopChanel)
	return nil
}

func (c *Scan) Start(ignore bool) error {
	c.stopped.Store(false)
	if c.ScanCall == nil {
		return errors.New("the ScanCall is nil")
	}
	go c.startScan()
	return nil
}

func (c *Scan) startScan() {
	c.initBlockNumber = big.NewInt(-1)

	for { //get the  init block number
		if c.stopped.Load() {
			return
		}
		t := c.ScanCall.StartBlockNumber()
		if t.Cmp(dao.GetErrorBigInt) != 0 {
			c.initBlockNumber.Set(t)
			break
		}
		select {
		case <-c.stopChanel:
			return
		case <-time.After(time.Second * IntervalScanCall):
		}
	}

	cur := big.NewInt(0)
	cur.Add(c.initBlockNumber, big.NewInt(1))
	eth := c.EthConnect
	max := big.NewInt(0)

ForMaxBlocks:
	for { //取最大区块号，遍历区块，取区块中的所有交易
		for { //get the max block number
			if c.stopped.Load() {
				break ForMaxBlocks
			}
			b, err := eth.EthClient.BlockByNumber(eth.Ctx, nil)
			if err == nil {
				max.Sub(b.Header().Number, big.NewInt(FinalHight))
				break
			}
			select { //等待后再重试
			case <-c.stopChanel:
				break ForMaxBlocks
			case <-time.After(time.Second * IntervalBlock):
			}
		}

		if cur.Cmp(max) >= 0 { //如果没有新的区块; 这个条件不要放到 ForBlocks， 因为那样会每一次完成都等 2s
			select { //等待后再取最大区块号
			case <-c.stopChanel:
				break ForMaxBlocks
			case <-time.After(time.Second * IntervalBlock):
			}
			continue ForMaxBlocks
		}

	ForBlocks:
		for { //从当前一直到max区块号。不能在这里++,及作判断，因为cur是多次使用的
			var b *types.Block
			var err error
			for {
				if c.stopped.Load() {
					break ForMaxBlocks
				}
				b, err = eth.EthClient.BlockByNumber(eth.Ctx, cur)
				if err == nil && b != nil {
					break
				}
				select { //等待后再重试
				case <-c.stopChanel:
					break ForMaxBlocks
				case <-time.After(time.Second * IntervalBlock):
				}
			}

			for !c.ScanCall.Block(b) { //一直到返回true, 可能没有准备好
				select { //等待后再重试
				case <-c.stopChanel:
					break ForMaxBlocks
				case <-time.After(time.Second * IntervalScanCall):
				}
			}

			for _, htx := range b.Transactions() { //处理区块中的交易
				if c.stopped.Load() {
					break ForMaxBlocks
				}
				for {
					tx, _, err := eth.EthClient.TransactionByHash(eth.Ctx, htx.Hash())
					receipt, err2 := eth.EthClient.TransactionReceipt(eth.Ctx, htx.Hash())
					if err == nil && err2 == nil && tx != nil && receipt != nil {
						for !c.ScanCall.Tx(b, tx, receipt) { //一直到返回true, 可能没有准备好
							select { //等待后再重试
							case <-c.stopChanel:
								break ForMaxBlocks
							case <-time.After(time.Second * IntervalScanCall):
							}
						}
						break
					}

					select { //等待后再重试
					case <-c.stopChanel:
						break ForMaxBlocks
					case <-time.After(time.Second * IntervalTx):
					}
				}
			}

			for !c.ScanCall.DoneBlock(b) { //一直到返回true, 可能没有准备好
				select { //等待后再重试
				case <-c.stopChanel:
					break ForMaxBlocks
				case <-time.After(time.Second * IntervalScanCall):
				}
			}

			if cur.Cmp(max) >= 0 { //检测是否已到最大区块
				break ForBlocks
			}
			cur.Add(cur, big.NewInt(1))
		}
	}
}

func newMonitor(conf interface{}) (d dot.Dot, err error) {
	scan := &Scan{stopChanel: make(chan bool)}
	d = scan
	return d, err
}

//ScanTypeLives
func ScanTypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeId: ScanTypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return newMonitor(conf)
			}},
		},
	}
	lives = append(lives, eth.EthConnectTypeLives()...)
	return lives
}

////ScanConfigTypeLives
//func ScanConfigTypeLives() *dot.ConfigTypeLives {
//	return &dot.ConfigTypeLives{
//		TypeIdConfig: ScanTypeId,
//		ConfigInfo:   &scanConfig{},
//	}
//}
