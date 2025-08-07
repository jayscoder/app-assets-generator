package image

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generator 图片资源生成器
type Generator struct {
	inputPath  string
	outputPath string
}

// NewGenerator 创建新的生成器
func NewGenerator(inputPath, outputPath string) *Generator {
	return &Generator{
		inputPath:  inputPath,
		outputPath: outputPath,
	}
}

// GenerateIOS 生成iOS图片资源
func (g *Generator) GenerateIOS() error {
	// 扫描输入目录的图片
	images, err := g.scanImages()
	if err != nil {
		return fmt.Errorf("扫描图片失败: %w", err)
	}
	
	// 生成iOS资源
	iosGen := NewIOSImageGenerator(g.inputPath, g.outputPath)
	return iosGen.Generate(images)
}

// GenerateAndroid 生成Android图片资源
func (g *Generator) GenerateAndroid() error {
	// 扫描输入目录的图片
	images, err := g.scanImages()
	if err != nil {
		return fmt.Errorf("扫描图片失败: %w", err)
	}
	
	// 生成Android资源
	androidGen := NewAndroidImageGenerator(g.inputPath, g.outputPath)
	return androidGen.Generate(images)
}

// ImageInfo 图片信息
type ImageInfo struct {
	Name      string   // 图片名称（不含扩展名和@2x等后缀）
	Files     []string // 相关文件列表
	Extension string   // 文件扩展名
	Has1x     bool     // 是否有1x图片
	Has2x     bool     // 是否有@2x图片
	Has3x     bool     // 是否有@3x图片
}

// scanImages 扫描图片目录
func (g *Generator) scanImages() (map[string]*ImageInfo, error) {
	images := make(map[string]*ImageInfo)
	
	// 遍历目录
	err := filepath.Walk(g.inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		
		// 获取文件名
		fileName := info.Name()
		
		// 检查是否为支持的图片格式
		ext := strings.ToLower(filepath.Ext(fileName))
		if !isSupportedImageFormat(ext) {
			return nil
		}
		
		// 解析图片名称和倍数
		baseName, scale := parseImageName(fileName)
		
		// 获取或创建ImageInfo
		if _, exists := images[baseName]; !exists {
			images[baseName] = &ImageInfo{
				Name:      baseName,
				Extension: ext,
				Files:     []string{},
			}
		}
		
		imageInfo := images[baseName]
		imageInfo.Files = append(imageInfo.Files, fileName)
		
		// 标记倍数
		switch scale {
		case "1x":
			imageInfo.Has1x = true
		case "2x":
			imageInfo.Has2x = true
		case "3x":
			imageInfo.Has3x = true
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return images, nil
}

// parseImageName 解析图片名称，返回基础名称和倍数
func parseImageName(fileName string) (baseName string, scale string) {
	// 去除扩展名
	nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	
	// 检查@2x、@3x后缀
	if strings.HasSuffix(nameWithoutExt, "@3x") {
		return strings.TrimSuffix(nameWithoutExt, "@3x"), "3x"
	}
	if strings.HasSuffix(nameWithoutExt, "@2x") {
		return strings.TrimSuffix(nameWithoutExt, "@2x"), "2x"
	}
	
	// 默认为1x
	return nameWithoutExt, "1x"
}

// isSupportedImageFormat 检查是否为支持的图片格式
func isSupportedImageFormat(ext string) bool {
	supportedFormats := []string{".png", ".jpg", ".jpeg", ".svg", ".pdf"}
	for _, format := range supportedFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	
	// 读取源文件
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	// 写入目标文件
	return os.WriteFile(dst, data, 0644)
}