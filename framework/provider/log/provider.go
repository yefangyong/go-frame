package log

import (
	"io"
	"strings"

	"github.com/yefangyong/go-frame/framework/provider/log/formatter"

	"github.com/yefangyong/go-frame/framework/provider/log/services"

	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/contract"
)

type HadeLogServiceProvider struct {
	// 日志驱动
	Driver string

	// 日志级别
	Level contract.LogLevel

	// 日志输出格式方法
	Formatter contract.Formatter

	// 日志上下文获取信息的函数
	CtxFielder contract.CtxFielder

	// 日志输出信息
	Output io.Writer
}

func (h *HadeLogServiceProvider) Register(container framework.Container) framework.NewInstance {
	if h.Driver == "" {
		configServicePro, err := container.Make(contract.ConfigKey)
		if err != nil {
			// 默认使用console
			return services.NewHadeConsoleLog
		}
		configService := configServicePro.(contract.Config)
		h.Driver = strings.ToLower(configService.GetString("log.Driver"))
	}
	switch h.Driver {
	case "single":
		return services.NewHadeSingleLog
	case "rotate":
		return services.NewHadeRotateLog
	case "custom":
		return services.NewHadeCustomLog
	case "console":
		return services.NewHadeConsoleLog
	default:
		return services.NewHadeConsoleLog
	}
}

func (h *HadeLogServiceProvider) Boot(container framework.Container) error {
	return nil
}

func (h *HadeLogServiceProvider) IsDefer() bool {
	return false
}

func (h *HadeLogServiceProvider) Params(container framework.Container) []interface{} {
	// 获取 configService
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	if h.Formatter == nil {
		h.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				h.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				h.Formatter = formatter.TextFormatter
			}
		}
	}

	if h.Level == contract.UnknownLevel {
		h.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			h.Level = GetLevel(configService.GetString("log.level"))
		}
	}
	// 定义五个参数
	return []interface{}{container, h.Level, h.CtxFielder, h.Formatter, h.Output}
}

func (h *HadeLogServiceProvider) Name() string {
	return contract.LogKey
}

func GetLevel(level string) contract.LogLevel {
	switch strings.ToLower(level) {
	case "panic":
		return contract.PanicLevel
	case "info":
		return contract.InfoLevel
	case "warn":
		return contract.WarnLevel
	case "error":
		return contract.ErrorLevel
	case "fatal":
		return contract.FatalLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknownLevel
}
