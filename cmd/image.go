package cmd

import (
	"app-assets-generator/pkg/image"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	imageInput    string
	imageOutput   string
	imagePlatform string
)

// imageCmd 图片生成命令
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "生成图片资源文件",
	Long:  `处理图片资源并生成iOS和Android平台的图片资源文件`,
	Example: `  # iOS平台
  app-assets-generator image --input icons/ --output output/ios --platform ios
  
  # Android平台
  app-assets-generator image --input icons/ --output output/android --platform android
  
  # 同时生成两个平台
  app-assets-generator image --input icons/ --output output/ --platform all`,
	Run: runImageCommand,
}

func init() {
	// 注册命令
	rootCmd.AddCommand(imageCmd)
	
	// 添加flag
	imageCmd.Flags().StringVarP(&imageInput, "input", "i", "", "输入的图片目录路径 (必需)")
	imageCmd.Flags().StringVarP(&imageOutput, "output", "o", "", "输出目录路径 (必需)")
	imageCmd.Flags().StringVarP(&imagePlatform, "platform", "p", "all", "目标平台 (ios/android/all)")
	
	// 标记必需的flag
	imageCmd.MarkFlagRequired("input")
	imageCmd.MarkFlagRequired("output")
}

func runImageCommand(cmd *cobra.Command, args []string) {
	// 验证输入目录是否存在
	if info, err := os.Stat(imageInput); os.IsNotExist(err) {
		exitWithError("输入目录不存在: %s", imageInput)
	} else if !info.IsDir() {
		exitWithError("输入路径不是目录: %s", imageInput)
	}
	
	// 验证平台参数
	if imagePlatform != "ios" && imagePlatform != "android" && imagePlatform != "all" {
		exitWithError("无效的平台参数: %s (必须是 ios/android/all)", imagePlatform)
	}
	
	// 创建生成器
	generator := image.NewGenerator(imageInput, imageOutput)
	
	// 根据平台生成资源
	var err error
	switch imagePlatform {
	case "ios":
		fmt.Println("正在生成iOS图片资源...")
		err = generator.GenerateIOS()
	case "android":
		fmt.Println("正在生成Android图片资源...")
		err = generator.GenerateAndroid()
	case "all":
		fmt.Println("正在生成iOS图片资源...")
		if err = generator.GenerateIOS(); err != nil {
			exitWithError("生成iOS资源失败: %v", err)
		}
		fmt.Println("正在生成Android图片资源...")
		err = generator.GenerateAndroid()
	}
	
	if err != nil {
		exitWithError("生成失败: %v", err)
	}
	
	fmt.Printf("✅ 图片资源生成成功！输出目录: %s\n", imageOutput)
}