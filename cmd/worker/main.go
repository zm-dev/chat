package main

import (
	"flag"
	"github.com/zm-dev/chat_v2/queue/subscribe"
	"github.com/zm-dev/chat_v2/server"
	"go.uber.org/zap"
)

var (
	h bool
	c string
)

func init() {
	flag.BoolVar(&h, "h", false, "the help")
	flag.StringVar(&c, "c", "config/config.yml", "set configuration `file`")
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return

	}
	svr := server.SetupServer(c)
	svr.Logger.Info("start queue", zap.Int("queue goroutine num", svr.Conf.QueueNum))
	subscribe.StartSubQueue(svr)
}
