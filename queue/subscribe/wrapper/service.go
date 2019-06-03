package wrapper

import (
	"context"
	"github.com/zm-dev/chat/pkg/pubsub"
	"github.com/zm-dev/chat/service"
)

type Service struct {
	sub     pubsub.SubQueue
	service service.Service
}

func (g *Service) Channel() string {
	return g.sub.Channel()
}

func (g *Service) Process(ctx context.Context, message string) {
	g.sub.Process(service.NewContext(ctx, g.service), message)
}

func NewService(sub pubsub.SubQueue, service service.Service) pubsub.SubQueue {
	return &Service{sub: sub, service: service}
}
