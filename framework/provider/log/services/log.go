package services

import (
	"context"
	"io"
	pkgLog "log"
	"time"

	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
	"github.com/yefangyong/go-frame/framework/provider/log/formatter"
)

// 日志的通用实例
type HadeLog struct {
	container  framework.Container
	level      contract.LogLevel
	ctxFielder contract.CtxFielder
	formatter  contract.Formatter
	output     io.Writer
}

func (h *HadeLog) logf(level contract.LogLevel, ctx context.Context, msg string, field map[string]interface{}) error {
	// 先判断日志级别
	if !h.IsLevelEnable(level) {
		return nil
	}
	fs := field
	// 使用 ctxFielder 获取 context 中的信息
	if h.ctxFielder != nil {
		t := h.ctxFielder(ctx)
		for k, v := range t {
			fs[k] = v
		}
	}

	// 将日志信息根据 formatter 格式化为字符串
	if h.formatter == nil {
		h.formatter = formatter.TextFormatter
	}
	ct, err := h.formatter(level, time.Now(), msg, field)
	if err != nil {
		return err
	}

	// 如果是panic级别，则使用日志进行panic
	if level == contract.PanicLevel {
		pkgLog.Panicln(string(ct))
		return nil
	}

	// 通过 output 进行输出
	_, _ = h.output.Write(ct)
	_, _ = h.output.Write([]byte("\r\n"))
	return nil
}

func (h *HadeLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	h.logf(contract.PanicLevel, ctx, msg, fields)
}

func (h *HadeLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	h.logf(contract.FatalLevel, ctx, msg, fields)
}

func (h *HadeLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	h.logf(contract.ErrorLevel, ctx, msg, fields)
}

func (h *HadeLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	h.logf(contract.WarnLevel, ctx, msg, fields)
}

func (h *HadeLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	h.logf(contract.InfoLevel, ctx, msg, fields)
}

func (h *HadeLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	h.logf(contract.DebugLevel, ctx, msg, fields)
}

func (h *HadeLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	h.logf(contract.TraceLevel, ctx, msg, fields)
}

func (h *HadeLog) SetLevel(level contract.LogLevel) {
	h.level = level
}

func (h *HadeLog) SetCtxFielder(handler contract.CtxFielder) {
	h.ctxFielder = handler
}

func (h *HadeLog) SetFormatter(formatter contract.Formatter) {
	h.formatter = formatter
}

func (h *HadeLog) SetOutput(out io.Writer) {
	h.output = out
}

// 判断这个日志级别是否可以打印
func (h *HadeLog) IsLevelEnable(level contract.LogLevel) bool {
	return h.level <= level
}
