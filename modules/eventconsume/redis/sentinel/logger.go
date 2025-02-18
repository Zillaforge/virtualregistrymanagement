package sentinel

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func init() {
	redis.SetLogger(&redisLogger{})
}

type redisLogger struct {
}

func (l *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	zap.L().Info(fmt.Sprintf(format, v...))
}
