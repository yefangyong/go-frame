package contract

import "time"

const ConfigKey = "hade:config"

type Config interface {
	IsExist(key string) bool

	Get(key string) interface{}

	GetString(key string) string

	GetInt(key string) int

	GetBool(key string) bool

	GetFloat64(key string) float64

	GetTime(key string) time.Time

	GetStringSlice(key string) []string

	GetIntSlice(key string) []int

	GetStringMap(key string) map[string]interface{}

	GetStringMapString(key string) map[string]string

	GetStringMapStringSlice(key string) map[string][]string

	Load(key string, val interface{}) error
}
