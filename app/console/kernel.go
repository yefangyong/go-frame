package console

import (
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

func AddAppCommand(rootCmd *cobra.Command) {

}
