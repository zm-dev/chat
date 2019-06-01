package wrapper

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat_v2/pkg/pubsub"
	"github.com/zm-dev/chat_v2/store/db_store"
)

type DB struct {
	sub pubsub.SubQueue
	db  *gorm.DB
}

func (g *DB) Channel() string {
	return g.sub.Channel()
}

func (g *DB) Process(ctx context.Context, message string) {
	g.sub.Process(db_store.NewDBContext(ctx, g.db), message)
}

func NewDB(sub pubsub.SubQueue, db *gorm.DB) pubsub.SubQueue {
	return &DB{sub: sub, db: db}
}
