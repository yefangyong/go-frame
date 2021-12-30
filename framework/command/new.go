package command

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/google/go-github/v39/github"

	"github.com/AlecAivazis/survey/v2"
	"github.com/yefangyong/go-frame/framework/cobra"
	"github.com/yefangyong/go-frame/framework/util"
)

func initNewCommand() *cobra.Command {
	return newCommand
}

var newCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个新的应用",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentPath := util.GetExecDirectory()

		var name string
		var folder string
		var mod string
		var version string
		var release *github.RepositoryRelease
		{
			prompt := &survey.Input{
				Message: "请输入目录名称：",
			}

			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}

			folder = filepath.Join(currentPath, name)
			if util.Exists(folder) {
				isForce := false
				prompt2 := &survey.Confirm{
					Message: "目录" + folder + "已经存在，是否删除重新创建（确认后立即执行删除操作！）",
					Default: false,
				}
				err := survey.AskOne(prompt2, &isForce)
				if err != nil {
					return err
				}

				if isForce {
					if err := os.RemoveAll(folder); err != nil {
						return err
					}
				} else {
					fmt.Println("目录已经存在，创建应用失败")
					return nil
				}
			}
		}
		fmt.Println(folder)
		{
			prompt := &survey.Input{
				Message: "请输入模块名称（go.mod中的module，默认为文件夹名称）:",
			}
			err := survey.AskOne(prompt, &mod)
			if err != nil {
				return err
			}
			if mod == "" {
				mod = name
			}
		}
		{
			// 获取hade的版本
			client := github.NewClient(nil)
			prompt := &survey.Input{
				Message: "请输入版本名称（参考 https://github.com/yefangyong/go-frame/releases,默认为最新版本号）",
			}
			err := survey.AskOne(prompt, &version)
			if err != nil {
				return err
			}
			if version != "" {
				release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "yefangyong", "go-frame", version)
				if err != nil || release == nil {
					fmt.Println("版本不存在，创建应用失败，请参考 https://github.com/gohade/hade/releases")
					return nil
				}
			}
			if version == "" {
				release, _, err = client.Repositories.GetLatestRelease(context.Background(), "yefangyong", "go-frame")
				if err != nil || release == nil {
					fmt.Println("版本不存在，创建应用失败，请参考 https://github.com/gohade/hade/releases")
					return nil
				}
				version = release.GetTagName()
			}

			fmt.Println("=================================")
			fmt.Println("开始创建应用操作")
			fmt.Println("创建目录：", folder)
			fmt.Println("应用名称：", mod)
			fmt.Println("hade框架版本：", release.GetTagName())

			templateFolder := filepath.Join(currentPath, "template-goFrame-"+version+"-"+cast.ToString(time.Now().Unix()))
			os.Mkdir(templateFolder, os.ModePerm)
			fmt.Println("创建临时目录", templateFolder)

			// 拷贝 template 项目
			url := release.GetZipballURL()
			err = util.DownLoadFile(filepath.Join(templateFolder, "template.zip"), url)
			if err != nil {
				return err
			}

			fmt.Println("下载zip包到template.zip")
			_, err = util.Unzip(filepath.Join(templateFolder, "template.zip"), templateFolder)
			if err != nil {
				return err
			}

			// 获取 folder 下的 goFrame-xxx 相关解压目录
			fInfos, err := ioutil.ReadDir(templateFolder)
			if err != nil {
				return err
			}
			for _, fInfo := range fInfos {
				// 找到解压后的文件夹
				if fInfo.IsDir() && strings.Contains(fInfo.Name(), "yefangyong-go-frame-") {
					if err := os.Rename(filepath.Join(templateFolder, fInfo.Name()), folder); err != nil {
						return err
					}
				}
			}
			fmt.Println("解压zip包")

			if err := os.RemoveAll(templateFolder); err != nil {
				return err
			}
			fmt.Println("删除临时文件夹", templateFolder)

			os.RemoveAll(path.Join(folder, ".git"))
			fmt.Println("删除.git目录")

			// 删除 framework 目录
			os.RemoveAll(path.Join(folder, "framework"))
			fmt.Println("删除 framework 目录")

			filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				c, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				if path == filepath.Join(folder, "go.mod") {
					fmt.Println("更新文件：" + path)
					c = bytes.ReplaceAll(c, []byte("module github.com/yefangyong/go-frame"), []byte("module "+mod))
					c = bytes.ReplaceAll(c, []byte("require ("), []byte("require (\n\tgithub.com/yefangyong/go-frame "+version))
					err = ioutil.WriteFile(path, c, 0644)
					if err != nil {
						return err
					}
					return nil
				}

				isContain := bytes.Contains(c, []byte("github.com/yefangyong/go-frame/app"))
				if isContain {
					fmt.Println("更新文件：" + path)
					c = bytes.ReplaceAll(c, []byte("github.com/yefangyong/go-frame/app"), []byte(mod+"/app"))
					err = ioutil.WriteFile(path, c, 0644)
					if err != nil {
						return err
					}
				}
				return nil
			})
		}
		fmt.Println("创建应用结束")
		fmt.Println("目录：", folder)
		fmt.Println("===========================================")
		return nil
	},
}
