package command

import "github.com/yefangyong/go-frame/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(initAppCommand())
}
