package gogo

import (
	"github.com/spf13/cobra"
	"github.com/wzyjerry/windranger/internal/command"
	"github.com/wzyjerry/windranger/internal/generator/gogo"
	"github.com/wzyjerry/windranger/internal/parser"
)

// Gogo 根据配置文件生成go文件
func Gogo() *cobra.Command {
	var cfg command.Config
	cmd := &cobra.Command{
		Use:   "gogo [flags] profile",
		Short: "根据配置文件生成go文件",
		Example: command.Examples(
			"windranger gogo model --out model",
			"windranger gogo http://example.com/model.git --out model",
		),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p := parser.NewParser()
			p.AddYamlPath(args[0])
			packages, errs := p.Parse()
			if errs != nil {
				panic(errs[0])
			}
			if err := gogo.Generate(packages, cfg.Out); err != nil {
				panic(err)
			}
		},
	}
	// 生成根目录
	cmd.Flags().StringVar(&cfg.Out, "out", ".", "生成根目录")
	return cmd
}
