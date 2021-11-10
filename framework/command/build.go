package command

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/yefangyong/go-frame/framework/cobra"
)

func initBuildCommand() *cobra.Command {
	buildCommand.AddCommand(buildFrontendCommand)
	buildCommand.AddCommand(buildBackendCommand)
	buildCommand.AddCommand(buildSelfCommand)
	buildCommand.AddCommand(buildAllCommand)
	return buildCommand
}

var buildBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "使用 Go 编译后端",
	RunE: func(c *cobra.Command, args []string) error {
		return buildSelfCommand.RunE(c, args)
	},
}

var buildAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时编译前后端",
	RunE: func(c *cobra.Command, args []string) error {
		err := buildFrontendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		err = buildBackendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		return nil
	},
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "编译相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var buildSelfCommand = &cobra.Command{
	Use:   "self",
	Short: "编译 hade 命令",
	RunE: func(c *cobra.Command, args []string) error {
		// 获取 path 路径下面的 npm 命令
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalf("请安装 go 在你的PATH路径下")
		}
		// 执行 build
		cmd := exec.Command(path, "build", "-o", "hade", "./")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("======== 后端编译失败 ========")
			fmt.Println(string(out))
			fmt.Println("======== 后端编译失败 ========")
			return err
		}
		fmt.Println("go build success, please run ./hade direct")
		fmt.Println("======== 后端编译成功 ========")
		return nil
	},
}

var buildFrontendCommand = &cobra.Command{
	Use:   "fronted",
	Short: "使用 npm 编译前端",
	RunE: func(c *cobra.Command, args []string) error {
		// 获取 path 路径下面的 npm 命令
		path, err := exec.LookPath("npm")
		if err != nil {
			log.Fatalf("请安装npm在你的PATH路径下")
		}

		// 执行 npm run build
		cmd := exec.Command(path, "run", "build")

		// 将输出保存在 out 中
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("======== 前端编译失败 ========")
			fmt.Println(string(out))
			fmt.Println("======== 前端编译失败 ========")
			return err
		}

		// 打印输出
		fmt.Print(string(out))
		fmt.Println("======== 前端编译成功 ========")
		return nil
	},
}
