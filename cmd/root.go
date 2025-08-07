package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "app-assets-generator",
	Short: "生成Android和iOS应用资源文件的命令行工具",
	Long: `App Assets Generator 是一个用于生成Android和iOS应用资源文件的命令行工具。
	
支持功能：
- 从YAML配置文件批量生成颜色资源
- 自动处理多分辨率图片资源
- 支持Light/Dark主题配置
- 同时支持iOS和Android平台`,
	Version: version,
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 设置版本输出模板
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s\n" .Version}}`)
	
	// 初始化时可以在这里添加全局flag
	// rootCmd.PersistentFlags().StringP("config", "c", "", "配置文件路径")
}

// 退出并显示错误
func exitWithError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "错误: "+msg+"\n", args...)
	os.Exit(1)
}