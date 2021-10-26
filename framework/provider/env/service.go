package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/yefangyong/go-frame/framework/contract"
)

type HadeEnv struct {
	folder string            // 代表.env所在的目录
	maps   map[string]string // 代表所有的环境变量
}

func (h *HadeEnv) AppEnv() string {
	return h.Get("APP_ENV")
}

func (h *HadeEnv) IsExist(key string) bool {
	_, ok := h.maps[key]
	return ok
}

func (h *HadeEnv) Get(key string) string {
	if val, ok := h.maps[key]; ok {
		return val
	}
	return ""
}

func (h *HadeEnv) All() map[string]string {
	return h.maps
}

func NewHadeEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewHadeEnv params error")
	}

	// 读取 folder 文件
	folder := params[0].(string)

	// 实例化
	hadeEnv := &HadeEnv{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDevelopment}, // 默认为开发环境
	}

	// 解析 folder/.env 文件
	file := path.Join(folder, ".env")

	fi, err := os.Open(file)

	if err == nil {
		defer fi.Close()

		// 读取文件
		br := bufio.NewReader(fi)
		for {
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}

			//按照等号进行解析
			s := bytes.SplitN(line, []byte{'='}, 2)

			// 如果不符合规范，则过滤
			if len(s) < 2 {
				continue
			}

			// 保存 map
			key := string(s[0])
			val := string(s[1])
			hadeEnv.maps[key] = val
		}
	}

	// 获取当前程序的环境变量，并且覆盖 .env 文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		hadeEnv.maps[pair[0]] = pair[1]
	}
	return hadeEnv, nil
}
