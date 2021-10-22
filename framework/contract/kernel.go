package contract

import "net/http"

const KernelKey = "hade:kernel"

// 提供框架最核心的结构
type Kernel interface {
	HttpEngine() http.Handler
}
