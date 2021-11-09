package formatter

import "github.com/yefangyong/go-frame/framework/contract"

func Prefix(level contract.LogLevel) string {
	prefix := ""
	switch level {
	case contract.DebugLevel:
		prefix = "[Debug]"
	case contract.ErrorLevel:
		prefix = "[Error]"
	case contract.FatalLevel:
		prefix = "[Fatal]"
	case contract.InfoLevel:
		prefix = "[Info]"
	case contract.PanicLevel:
		prefix = "[Panic]"
	case contract.TraceLevel:
		prefix = "[Trace]"
	case contract.WarnLevel:
		prefix = "[Warn]"
	}
	return prefix
}
