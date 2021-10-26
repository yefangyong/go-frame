package contract

const (
	// 开发环境
	EnvDevelopment = "development"

	// 测试环境
	EnvTesting = "testing"

	// 生产环境
	EnvProduction = "production"

	EnvKey = "hade:env"
)

type Env interface {
	// AppEnv 获取当前app的环境，建议分为development,testing,production
	AppEnv() string

	//Get() 获取某个环境变量，如果没有设置，则返回""
	Get(key string) string

	// IsExist() 判断环境变量是否设置
	IsExist(key string) bool

	// 获取所有环境变量的值，.env和运行环境变量融合之后的结果
	All() map[string]string
}
