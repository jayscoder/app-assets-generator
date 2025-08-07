package cmd

import (
	"app-assets-generator/pkg/color"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	colorInput    string
	colorOutput   string
	colorPlatform string
)

// colorCmd 颜色生成命令
var colorCmd = &cobra.Command{
	Use:   "color",
	Short: "生成颜色资源文件",
	Long:  `从YAML配置文件生成iOS和Android平台的颜色资源文件`,
	Example: `  # iOS平台
  app-assets-generator color --input colors.yaml --output output/ios --platform ios
  
  # Android平台
  app-assets-generator color --input colors.yaml --output output/android --platform android
  
  # 同时生成两个平台
  app-assets-generator color --input colors.yaml --output output/ --platform all`,
	Run: runColorCommand,
}

func init() {
	// 注册命令
	rootCmd.AddCommand(colorCmd)
	
	// 添加flag
	colorCmd.Flags().StringVarP(&colorInput, "input", "i", "", "输入的YAML配置文件路径 (必需)")
	colorCmd.Flags().StringVarP(&colorOutput, "output", "o", "", "输出目录路径 (必需)")
	colorCmd.Flags().StringVarP(&colorPlatform, "platform", "p", "all", "目标平台 (ios/android/all)")
	
	// 标记必需的flag
	colorCmd.MarkFlagRequired("input")
	colorCmd.MarkFlagRequired("output")
}

func runColorCommand(cmd *cobra.Command, args []string) {
	// 验证输入文件是否存在
	if _, err := os.Stat(colorInput); os.IsNotExist(err) {
		exitWithError("输入文件不存在: %s", colorInput)
	}
	
	// 验证平台参数
	if colorPlatform != "ios" && colorPlatform != "android" && colorPlatform != "all" {
		exitWithError("无效的平台参数: %s (必须是 ios/android/all)", colorPlatform)
	}
	
	// 创建生成器
	generator := color.NewGenerator(colorInput, colorOutput)
	
	// 根据平台生成资源
	var err error
	switch colorPlatform {
	case "ios":
		fmt.Println("正在生成iOS颜色资源...")
		err = generator.GenerateIOS()
	case "android":
		fmt.Println("正在生成Android颜色资源...")
		err = generator.GenerateAndroid()
	case "all":
		fmt.Println("正在生成iOS颜色资源...")
		if err = generator.GenerateIOS(); err != nil {
			exitWithError("生成iOS资源失败: %v", err)
		}
		fmt.Println("正在生成Android颜色资源...")
		err = generator.GenerateAndroid()
	}
	
	if err != nil {
		exitWithError("生成失败: %v", err)
	}
	
	fmt.Printf("✅ 颜色资源生成成功！输出目录: %s\n", colorOutput)
}