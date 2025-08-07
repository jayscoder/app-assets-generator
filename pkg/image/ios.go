package image

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// IOSImageGenerator iOS图片资源生成器
type IOSImageGenerator struct {
	inputPath  string
	outputPath string
}

// NewIOSImageGenerator 创建iOS图片生成器
func NewIOSImageGenerator(inputPath, outputPath string) *IOSImageGenerator {
	return &IOSImageGenerator{
		inputPath:  inputPath,
		outputPath: outputPath,
	}
}

// iOSImageSet iOS图片集结构
type iOSImageSet struct {
	Images []iOSImage `json:"images"`
	Info   iOSInfo    `json:"info"`
}

// iOSImage iOS图片定义
type iOSImage struct {
	Filename string `json:"filename,omitempty"`
	Idiom    string `json:"idiom"`
	Scale    string `json:"scale"`
}

// iOSInfo iOS信息
type iOSInfo struct {
	Author  string `json:"author"`
	Version int    `json:"version"`
}

// Generate 生成iOS图片资源
func (g *IOSImageGenerator) Generate(images map[string]*ImageInfo) error {
	// 创建输出目录
	if err := os.MkdirAll(g.outputPath, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	
	// 为每个图片生成imageset
	for _, imageInfo := range images {
		if err := g.generateImageSet(g.outputPath, imageInfo); err != nil {
			return fmt.Errorf("生成图片 %s 失败: %w", imageInfo.Name, err)
		}
	}
	
	return nil
}

// generateImageSet 生成单个图片集
func (g *IOSImageGenerator) generateImageSet(outputPath string, imageInfo *ImageInfo) error {
	// 创建imageset目录
	imagesetPath := filepath.Join(outputPath, imageInfo.Name+".imageset")
	if err := os.MkdirAll(imagesetPath, 0755); err != nil {
		return fmt.Errorf("创建imageset目录失败: %w", err)
	}
	
	// 构建图片集数据
	imageSet := g.buildImageSet(imageInfo)
	
	// 复制图片文件
	for _, fileName := range imageInfo.Files {
		src := filepath.Join(g.inputPath, fileName)
		
		// 确定目标文件名
		dstFileName := g.getIOSFileName(imageInfo, fileName)
		dst := filepath.Join(imagesetPath, dstFileName)
		
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("复制图片文件失败: %w", err)
		}
	}
	
	// 生成Contents.json
	contentsPath := filepath.Join(imagesetPath, "Contents.json")
	file, err := os.Create(contentsPath)
	if err != nil {
		return fmt.Errorf("创建Contents.json失败: %w", err)
	}
	defer file.Close()
	
	// 写入JSON（格式化输出）
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(imageSet); err != nil {
		return fmt.Errorf("写入JSON失败: %w", err)
	}
	
	return nil
}

// buildImageSet 构建iOS图片集数据
func (g *IOSImageGenerator) buildImageSet(imageInfo *ImageInfo) iOSImageSet {
	imageSet := iOSImageSet{
		Images: []iOSImage{},
		Info: iOSInfo{
			Author:  "xcode",
			Version: 1,
		},
	}
	
	// 添加1x图片
	if imageInfo.Has1x {
		fileName := imageInfo.Name + imageInfo.Extension
		imageSet.Images = append(imageSet.Images, iOSImage{
			Filename: fileName,
			Idiom:    "universal",
			Scale:    "1x",
		})
	} else {
		// 如果没有1x图片，添加空占位
		imageSet.Images = append(imageSet.Images, iOSImage{
			Idiom: "universal",
			Scale: "1x",
		})
	}
	
	// 添加2x图片
	if imageInfo.Has2x {
		fileName := imageInfo.Name + "@2x" + imageInfo.Extension
		imageSet.Images = append(imageSet.Images, iOSImage{
			Filename: fileName,
			Idiom:    "universal",
			Scale:    "2x",
		})
	} else {
		// 如果没有2x图片，添加空占位
		imageSet.Images = append(imageSet.Images, iOSImage{
			Idiom: "universal",
			Scale: "2x",
		})
	}
	
	// 添加3x图片
	if imageInfo.Has3x {
		fileName := imageInfo.Name + "@3x" + imageInfo.Extension
		imageSet.Images = append(imageSet.Images, iOSImage{
			Filename: fileName,
			Idiom:    "universal",
			Scale:    "3x",
		})
	} else {
		// 如果没有3x图片，添加空占位
		imageSet.Images = append(imageSet.Images, iOSImage{
			Idiom: "universal",
			Scale: "3x",
		})
	}
	
	return imageSet
}

// getIOSFileName 获取iOS的文件名
func (g *IOSImageGenerator) getIOSFileName(imageInfo *ImageInfo, originalFileName string) string {
	// 保持原始文件名，包括@2x、@3x后缀
	return originalFileName
}