package api_service

import (
	"github.com/scryinfo/citychain/dots/data/model"
	"time"
)

type MustPurchases struct {
	pps     chan model.Purchase
	stopped chan bool
	call    func(*model.Purchase) bool
}

func (c *MustPurchases) Init(pps []model.Purchase, interval int64, call func(purchase *model.Purchase) bool) {

	c.pps = make(chan model.Purchase, 1000)
	c.stopped = make(chan bool)
	c.call = call
	go func() {
		for i, _ := range pps {
			for !call(&pps[i]) {
				select {
				case <-c.stopped:
					return
				case <-time.After(time.Second * time.Duration(interval)):
				}
			}
		}

		for {
			select {
			case <-c.stopped:
				return
			case pp := <-c.pps:
				for !call(&pp) {
					select {
					case <-c.stopped:
						return
					case <-time.After(time.Second * time.Duration(interval)):
					}
				}
			}
		}
	}()
}

//将数据放入缓冲channel， 如果channel已满，函数调用会阻塞
func (c *MustPurchases) Append(pp *model.Purchase) {
	c.pps <- *pp
}

func (c *MustPurchases) Stop() {
	if c.stopped != nil {
		close(c.stopped)
	}
}
