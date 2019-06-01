package main

import (
	"flag"
	"github.com/rs/cors"
	"github.com/zm-dev/chat_v2/handler"
	"github.com/zm-dev/chat_v2/server"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	svr.Logger.Info("listen", zap.String("addr", svr.Conf.ServerAddr))
	// cors 跨域用
	log.Fatal(http.ListenAndServe(svr.Conf.ServerAddr, cors.New(cors.Options{
		AllowedOrigins:   []string{"http://10.102.9.5", "http://10.102.9.5:81"},
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PUT", "HEAD"},
		AllowCredentials: true,
	}).Handler(handler.CreateHTTPHandler(svr))))
}
