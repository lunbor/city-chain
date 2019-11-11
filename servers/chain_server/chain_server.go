package main

import (
	"github.com/scryinfo/citychain/dots"
	"github.com/scryinfo/citychain/dots/api_service"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/line"
	"github.com/scryinfo/scryg/sutils/ssignal"
	"go.uber.org/zap"
	"os"
)

//在使用go build时默认生成的文件名是main所在的文件夹名字，所以命名为“chain_server”

func main() {
	l, err := line.BuildAndStartBy(&dot.Builder{
		Add: add,
	})
	if err != nil {
		dot.Logger().Errorln("", zap.Error(err))
		return
	}
	defer line.StopAndDestroy(l, true)
	dot.Logger().Infoln("dot ok")

	ssignal.WaitCtrlC(func(s os.Signal) bool {
		return false
	})
	dot.Logger().Infoln("dot will stop")
}

func add(l dot.Line) error {

	lives := dots.DoTxTypeLives() //todo
	lives = append(lives, api_service.TypeLives()...)

	err := l.PreAdd(lives...)

	return err
}
