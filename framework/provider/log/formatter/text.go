package formatter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/yefangyong/go-frame/framework/contract"
)

// 文本格式化
func TextFormatter(level contract.LogLevel, t time.Time, msg string, field map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	Separator := "\t"

	prefix := Prefix(level)
	bf.WriteString(prefix)
	bf.WriteString(Separator)

	// 输出时间
	ts := t.Format(time.RFC3339)
	bf.WriteString(ts)
	bf.WriteString(Separator)

	// 输出msg
	bf.WriteString("\"")
	bf.WriteString(msg)
	bf.WriteString("\"")
	bf.WriteString(Separator)

	// 输出 field
	bf.WriteString(fmt.Sprint(field))

	return bf.Bytes(), nil
}
