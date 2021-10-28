package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/spf13/cast"

	"github.com/fsnotify/fsnotify"

	"github.com/ghodss/yaml"

	"github.com/pkg/errors"

	"github.com/yefangyong/go-frame/framework"
)

type HadeConfig struct {
	container framework.Container    // 容器
	folder    string                 // 文件夹
	keyDelim  string                 // 路径的分隔符，默认为点
	lock      sync.RWMutex           // 配置文件的读写锁
	envMaps   map[string]string      // 所有的环境变量
	confMaps  map[string]interface{} //配置文件结构，以 key 为文件名
	confRaws  map[string][]byte      // 配置文件的原始信息
}

// 查找某个路径的配置项
func searchMap(maps map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return maps
	}
	// 判断是否有下个路径
	next, ok := maps[path[0]]
	if ok {
		// 判断这个路径是否为1
		if len(path) == 1 {
			return next
		}
		// 判断下一个路径的类型
		switch next.(type) {
		case map[interface{}]interface{}:
			// 如果是interface的map，使用cast进行value转换
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil
		}
	}
	return nil
}

// 通过path来获取某个配置项
func (conf *HadeConfig) find(key string) interface{} {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

func (conf *HadeConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

func (conf *HadeConfig) Get(key string) interface{} {
	return conf.find(key)
}

func (conf *HadeConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

func (conf *HadeConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

func (conf *HadeConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

func (conf *HadeConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

func (conf *HadeConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

func (conf *HadeConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

func (conf *HadeConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

func (conf *HadeConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

func (conf *HadeConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

func (conf *HadeConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

func (conf *HadeConfig) Load(key string, val interface{}) error {
	return mapstructure.Decode(conf.find(key), val)
}

func NewHadeConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	// 检查文件夹是否存在
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
	}

	// 实例化
	hadeConf := &HadeConfig{
		container: container,
		envMaps:   envMaps,
		lock:      sync.RWMutex{},
		keyDelim:  ".",
		folder:    envFolder,
		confMaps:  map[string]interface{}{},
		confRaws:  map[string][]byte{},
	}

	// 读取每一个文件
	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range files {
		fileName := file.Name()
		err := hadeConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// 监控文件夹文件，配置文件热更新
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watch.Add(envFolder)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			select {
			case ev := <-watch.Events:
				{
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件：", ev.Name)
						hadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件：", ev.Name)
						hadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件：", ev.Name)
						hadeConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error: ", err)
					return
				}
			}
		}
	}()
	return hadeConf, nil
}

// 删除某个配置文件
func (conf *HadeConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confMaps, name)
		delete(conf.confRaws, name)
	}
	return nil
}

// 读取某个配置文件
func (conf *HadeConfig) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	// 判断文件是否以为yaml或者yml作为后缀
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		// 读取文件的内容
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}

		// 直接针对文本做环境变量的替换
		bf = replace(bf, conf.envMaps)

		// 解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		conf.confRaws[name] = bf
		conf.confMaps[name] = c
	}
	return nil
}

func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	// 直接使用ReplaceAll替换，性能可能不是最优，但是配置文件加载，频率很低，问题不大
	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}
	return content
}
