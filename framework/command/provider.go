package command

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yefangyong/go-frame/framework/contract"
	"github.com/yefangyong/go-frame/framework/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jianfengye/collection"
	"github.com/yefangyong/go-frame/framework"

	"github.com/yefangyong/go-frame/framework/cobra"
)

// 初始化 provider 相关服务
func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerCreateCommand)
	providerCommand.AddCommand(providerListCommand)
	return providerCommand
}

// providerCommand 二级命令
var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供者相关命令",
	RunE: func(command *cobra.Command, args []string) error {
		if len(args) == 0 {
			command.Help()
		}
		return nil
	},
}

// 列出容器内所有的服务提供者
var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出容器内的所有服务",
	RunE: func(command *cobra.Command, args []string) error {
		container := command.GetContainer()
		hadeContainer := container.(*framework.HadeContainer)
		list := hadeContainer.ProviderList()
		for _, line := range list {
			println(line)
		}
		return nil
	},
}

var providerCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		fmt.Println("创建一个服务")
		var name, folder string
		{
			prompt := &survey.Input{
				Message: "请输入服务名称(服务凭证):",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入服务所在目录名称（默认：同服务名称）:",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}

		// 检查服务是否存在
		providers := container.(*framework.HadeContainer).ProviderList()
		providerColl := collection.NewStrCollection(providers)
		if providerColl.Contains(name) {
			fmt.Println("服务名称已经存在")
			return nil
		}

		if folder == "" {
			folder = name
		}

		app := container.MustMake(contract.AppKey).(contract.App)
		providerFolder := app.ProviderFolder()
		subFolders, err := util.SubDir(providerFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
			return nil
		}

		// 开始创建目录
		if err := os.Mkdir(filepath.Join(providerFolder, folder), 0700); err != nil {
			return err
		}

		funcs := template.FuncMap{"title": strings.Title}
		{
			//创建 contract.go
			file := filepath.Join(providerFolder, folder, "contract.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}

			// 使用 contractTmp 模板来初始化 template,并且让这个模板支持 title 方法，即支持{{.|title}}
			t := template.Must(template.New("contract").Funcs(funcs).Parse(contractTmp))
			// 将 name 传递到 template 中渲染，并且输出到 contract.go 中
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		{
			//创建 provider.go
			file := filepath.Join(providerFolder, folder, "provider.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}

			// 使用 providerTmp 模板来初始化 template,并且让这个模板支持 title 方法，即支持{{.|title}}
			t := template.Must(template.New("provider").Funcs(funcs).Parse(providerTmp))
			// 将 name 传递到 template 中渲染，并且输出到 provider.go 中
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		{
			//创建 service.go
			file := filepath.Join(providerFolder, folder, "service.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}

			// 使用 serviceTmp 模板来初始化 template,并且让这个模板支持 title 方法，即支持{{.|title}}
			t := template.Must(template.New("service").Funcs(funcs).Parse(serviceTmp))
			// 将 name 传递到 template 中渲染，并且输出到 service.go 中
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		fmt.Println("创建服务成功，文件夹地址：", filepath.Join(providerFolder, folder))
		fmt.Println("请不要忘记挂载新创建的服务")
		return nil
	},
}

var contractTmp = `package {{.}}

const {{.|title}}Key = "{{.}}"

type Service interface {
	// 请在这里定义你的方法
    Foo() string
}
`

var providerTmp = `package {{.}}

import (
  "github.com/yefangyong/go-frame/framework"
)

type {{.|title}}Provider struct {
  framework.ServiceProvider

  c framework.Container
}

func (sp *{{.|title}}Provider) Name() string {
	return {{.|title}}Key
}

func (sp *{{.|title}}Provider) Register(c framework.Container) framework.NewInstance {
	return New{{.|title}}Service
}

func (sp *{{.|title}}Provider) IsDefer() bool {
	return false
}

func (sp *{{.|title}}Provider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (sp *{{.|title}}Provider) Boot(c framework.Container) error {
	return nil
}
`

var serviceTmp = `package {{.}}

import "github.com/yefangyong/go-frame/framework"

type {{.|title}}Service struct {
	container framework.Container
}

func New{{.|title}}Service(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &{{.|title}}Service{container: container}, nil
}

func (s *{{.|title}}Service) Foo() string {
    return ""
}
`
