package orm

import (
	"context"
	"time"

	"github.com/yefangyong/go-frame/framework/contract"
	"gorm.io/gorm/logger"
)

type OrmLogger struct {
	logger contract.Log
}

func (o *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
	return o
}

func (o OrmLogger) Info(ctx context.Context, s string, i ...interface{}) {
	o.logger.Info(ctx, s, map[string]interface{}{
		"filed": i,
	})
}

func (o OrmLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	o.logger.Warn(ctx, s, map[string]interface{}{
		"filed": i,
	})
}

func (o OrmLogger) Error(ctx context.Context, s string, i ...interface{}) {
	o.logger.Error(ctx, s, map[string]interface{}{
		"filed": i,
	})
}

func (o OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	fields := map[string]interface{}{
		"begin": begin,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
		"err":   err,
	}
	s := "orm trace sql"
	o.logger.Trace(ctx, s, fields)
}

// 初始化OrmLogger，使用框架自带的 logger
func NewOrmLogger(loggerService contract.Log) *OrmLogger {
	return &OrmLogger{
		logger: loggerService,
	}
}
