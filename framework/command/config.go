package command

import (
	"fmt"

	"github.com/kr/pretty"

	"github.com/yefangyong/go-frame/framework/cobra"
	"github.com/yefangyong/go-frame/framework/contract"
)

// initConfigCommand 获取配置相关的命令
func initConfigCommand() *cobra.Command {
	configCommand.AddCommand(configGetCommand)
	return configCommand
}

// 二级命令
var configCommand = &cobra.Command{
	Use:   "config",
	Short: "获取配置相关的信息",
	RunE: func(command *cobra.Command, args []string) error {
		if len(args) == 0 {
			command.Help()
		}
		return nil
	},
}

// 获取配置命令
var configGetCommand = &cobra.Command{
	Use:   "get",
	Short: "获取某个配置信息",
	RunE: func(command *cobra.Command, args []string) error {
		container := command.GetContainer()
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		if len(args) != 1 {
			fmt.Println("参数错误")
			return nil
		}
		key := args[0]
		val := configService.Get(key)
		if val == nil {
			fmt.Println("配置路径 ", configService, " 不存在")
			return nil
		}
		fmt.Printf("%# v\n", pretty.Formatter(val))
		return nil
	},
}
