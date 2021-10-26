package env

import (
	"fmt"
	"testing"

	"github.com/yefangyong/go-frame/framework/contract"
	"github.com/yefangyong/go-frame/framework/provider/app"

	"github.com/yefangyong/go-frame/framework"
)

func TestHadeEnvProvider(t *testing.T) {

	c := framework.NewHadeContainer()
	sp := &app.HadeAppProvider{}

	_ = c.Bind(sp)
	//So(err, ShouldBeNil)

	sp2 := &HadeEnvProvider{}
	_ = c.Bind(sp2)
	//So(err, ShouldBeNil)

	envServ := c.MustMake(contract.EnvKey).(contract.Env)
	fmt.Println(envServ.All())
	//So(envServ.AppEnv(), ShouldEqual, "development")
	// So(envServ.Get("DB_HOST"), ShouldEqual, "127.0.0.1")
	// So(envServ.AppDebug(), ShouldBeTrue)
}
