package command

import "strings"

type (
	// Config 配置
	Config struct {
		// Out 生成根目录
		Out string
	}
)

// Examples 格式化多个示例用法
func Examples(values ...string) string {
	// 添加两个前导空格
	for i, value := range values {
		values[i] = "  " + value
	}
	return strings.Join(values, "\n")
}
