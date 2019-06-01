package subscribe

import (
	"context"
	"github.com/zm-dev/chat_v2/pkg/pubsub"
	"github.com/zm-dev/chat_v2/server"
)

func StartSubQueue(svr *server.Server) {
	ctx := context.Background()
	sub := pubsub.NewSub(svr.RedisClient, svr.Logger, svr.Conf.QueueNum)
	sub.RegisterSub()
	sub.Sub(ctx)
}
