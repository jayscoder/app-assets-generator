package image

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AndroidImageGenerator Android图片资源生成器
type AndroidImageGenerator struct {
	inputPath  string
	outputPath string
}

// NewAndroidImageGenerator 创建Android图片生成器
func NewAndroidImageGenerator(inputPath, outputPath string) *AndroidImageGenerator {
	return &AndroidImageGenerator{
		inputPath:  inputPath,
		outputPath: outputPath,
	}
}

// AndroidDensity Android屏幕密度配置
type AndroidDensity struct {
	Name      string  // 密度名称
	Scale     float64 // 相对于mdpi的缩放比例
	Directory string  // 目录名称
}

var androidDensities = []AndroidDensity{
	{Name: "mdpi", Scale: 1.0, Directory: "drawable-mdpi"},      // 1x
	{Name: "hdpi", Scale: 1.5, Directory: "drawable-hdpi"},      // 1.5x
	{Name: "xhdpi", Scale: 2.0, Directory: "drawable-xhdpi"},    // 2x
	{Name: "xxhdpi", Scale: 3.0, Directory: "drawable-xxhdpi"},  // 3x
	{Name: "xxxhdpi", Scale: 4.0, Directory: "drawable-xxxhdpi"}, // 4x
}

// Generate 生成Android图片资源
func (g *AndroidImageGenerator) Generate(images map[string]*ImageInfo) error {
	// 为每个图片生成Android资源
	for _, imageInfo := range images {
		if err := g.generateAndroidImage(imageInfo); err != nil {
			return fmt.Errorf("生成图片 %s 失败: %w", imageInfo.Name, err)
		}
	}
	
	return nil
}

// generateAndroidImage 生成单个Android图片资源
func (g *AndroidImageGenerator) generateAndroidImage(imageInfo *ImageInfo) error {
	// Android使用下划线命名，将连字符转换为下划线
	androidName := strings.ReplaceAll(imageInfo.Name, "-", "_")
	androidName = strings.ToLower(androidName) // Android资源名称通常使用小写
	
	// 根据可用的iOS图片决定如何分配到Android密度
	mapping := g.getAndroidMapping(imageInfo)
	
	// 复制图片到对应的drawable目录
	for density, sourceFile := range mapping {
		if sourceFile == "" {
			continue // 跳过没有对应源文件的密度
		}
		
		// 创建目标目录
		targetDir := filepath.Join(g.outputPath, density.Directory)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", targetDir, err)
		}
		
		// 复制文件
		src := filepath.Join(g.inputPath, sourceFile)
		dst := filepath.Join(targetDir, androidName+imageInfo.Extension)
		
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("复制图片文件失败: %w", err)
		}
	}
	
	return nil
}

// getAndroidMapping 获取iOS图片到Android密度的映射
func (g *AndroidImageGenerator) getAndroidMapping(imageInfo *ImageInfo) map[AndroidDensity]string {
	mapping := make(map[AndroidDensity]string)
	
	// 理想映射：
	// iOS 1x -> Android mdpi (1x)
	// iOS 2x -> Android xhdpi (2x)
	// iOS 3x -> Android xxhdpi (3x)
	
	// 如果有1x图片，用于mdpi
	if imageInfo.Has1x {
		mapping[androidDensities[0]] = imageInfo.Name + imageInfo.Extension // mdpi
	}
	
	// 如果有2x图片，用于hdpi和xhdpi
	if imageInfo.Has2x {
		fileName2x := imageInfo.Name + "@2x" + imageInfo.Extension
		if !imageInfo.Has1x {
			// 如果没有1x，2x也用于mdpi
			mapping[androidDensities[0]] = fileName2x // mdpi
		}
		mapping[androidDensities[1]] = fileName2x // hdpi
		mapping[androidDensities[2]] = fileName2x // xhdpi
	}
	
	// 如果有3x图片，用于xxhdpi和xxxhdpi
	if imageInfo.Has3x {
		fileName3x := imageInfo.Name + "@3x" + imageInfo.Extension
		mapping[androidDensities[3]] = fileName3x // xxhdpi
		mapping[androidDensities[4]] = fileName3x // xxxhdpi
		
		// 如果没有2x，3x也用于xhdpi
		if !imageInfo.Has2x {
			mapping[androidDensities[2]] = fileName3x // xhdpi
			// 如果连1x都没有，3x用于所有密度
			if !imageInfo.Has1x {
				mapping[androidDensities[0]] = fileName3x // mdpi
				mapping[androidDensities[1]] = fileName3x // hdpi
			}
		}
	}
	
	// 如果只有1x图片，用于所有密度
	if imageInfo.Has1x && !imageInfo.Has2x && !imageInfo.Has3x {
		fileName1x := imageInfo.Name + imageInfo.Extension
		for i := range androidDensities {
			mapping[androidDensities[i]] = fileName1x
		}
	}
	
	return mapping
}