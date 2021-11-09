package formatter

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/yefangyong/go-frame/framework/contract"
)

// json格式化日志信息
func JsonFormatter(level contract.LogLevel, t time.Time, msg string, field map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	field["msg"] = msg
	field["timestamp"] = t.Format(time.RFC3339)
	field["level"] = level
	c, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}
	bf.WriteString(string(c))
	return bf.Bytes(), nil
}
