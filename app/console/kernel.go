package console

import (
	"time"

	"github.com/yefangyong/go-frame/app/console/command/demo"
	"github.com/yefangyong/go-frame/framework"
	"github.com/yefangyong/go-frame/framework/cobra"
	"github.com/yefangyong/go-frame/framework/command"
)

func RunCommand(container framework.Container) error {
	// 根Command
	var rootCmd = &cobra.Command{
		Use:   "hade",
		Short: "hade 命令",
		Long:  "hade 框架提供的命令行工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	// 为根Command设置服务容器
	rootCmd.SetContainer(container)

	// 绑定框架的命令
	command.AddKernelCommands(rootCmd)

	// 绑定业务的命令
	AddAppCommand(rootCmd)

	// 执行RootCommand
	return rootCmd.Execute()
}

// 业务的相关命令
func AddAppCommand(rootCmd *cobra.Command) {

	rootCmd.AddCommand(demo.InitFoo())

	//rootCmd.AddCronCommand("* * * * * *", demo.Foo1Command)

	// 启动一个分布式任务调度，调度的服务名称为init_func_for_test，每个节点每5s调用一次Foo命令，抢占到了调度任务的节点将抢占锁持续挂载2s才释放
	rootCmd.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}
