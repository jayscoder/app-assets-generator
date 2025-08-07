package color

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AndroidGenerator Android颜色资源生成器
type AndroidGenerator struct {
	outputPath string
}

// NewAndroidGenerator 创建Android生成器
func NewAndroidGenerator(outputPath string) *AndroidGenerator {
	return &AndroidGenerator{
		outputPath: outputPath,
	}
}

// Generate 生成Android颜色资源
func (g *AndroidGenerator) Generate(colors map[string]*ColorDefinition) error {
	// 创建values目录
	valuesPath := filepath.Join(g.outputPath, "values")
	if err := os.MkdirAll(valuesPath, 0755); err != nil {
		return fmt.Errorf("创建values目录失败: %w", err)
	}
	
	// 创建values-night目录（深色主题）
	valuesNightPath := filepath.Join(g.outputPath, "values-night")
	if err := os.MkdirAll(valuesNightPath, 0755); err != nil {
		return fmt.Errorf("创建values-night目录失败: %w", err)
	}
	
	// 收集默认颜色和深色主题颜色
	defaultColors := make(map[string]string)
	nightColors := make(map[string]string)
	
	for name, color := range colors {
		// 跳过渐变色
		if color.IsGradient() {
			continue
		}
		
		// 获取默认/浅色主题颜色
		lightColor := color.GetLight()
		if lightColor.Hex != "" {
			defaultColors[name] = g.formatAndroidColor(lightColor)
		}
		
		// 获取深色主题颜色
		darkColor := color.GetDark()
		if darkColor.Hex != "" {
			// 只有当深色主题颜色与浅色不同时才添加
			if darkColor.Hex != lightColor.Hex || darkColor.Alpha != lightColor.Alpha {
				nightColors[name] = g.formatAndroidColor(darkColor)
			}
		}
	}
	
	// 生成默认colors.xml
	if err := g.generateColorsXML(valuesPath, defaultColors); err != nil {
		return fmt.Errorf("生成默认colors.xml失败: %w", err)
	}
	
	// 如果有深色主题颜色，生成values-night/colors.xml
	if len(nightColors) > 0 {
		if err := g.generateColorsXML(valuesNightPath, nightColors); err != nil {
			return fmt.Errorf("生成深色主题colors.xml失败: %w", err)
		}
	}
	
	return nil
}

// generateColorsXML 生成colors.xml文件
func (g *AndroidGenerator) generateColorsXML(dirPath string, colors map[string]string) error {
	filePath := filepath.Join(dirPath, "colors.xml")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建colors.xml失败: %w", err)
	}
	defer file.Close()
	
	// 写入XML头
	fmt.Fprintln(file, `<?xml version="1.0" encoding="utf-8"?>`)
	fmt.Fprintln(file, `<resources>`)
	
	// 按字母顺序排序并写入颜色
	sortedNames := g.sortColorNames(colors)
	for _, name := range sortedNames {
		colorValue := colors[name]
		// 将下划线转换为更符合Android命名规范（可选）
		androidName := name
		fmt.Fprintf(file, `    <color name="%s">%s</color>`+"\n", androidName, colorValue)
	}
	
	// 写入结束标签
	fmt.Fprintln(file, `</resources>`)
	
	return nil
}

// formatAndroidColor 格式化Android颜色值
func (g *AndroidGenerator) formatAndroidColor(color ColorValue) string {
	// Android颜色格式: #AARRGGBB 或 #RRGGBB
	hex := color.Hex
	
	// 如果透明度不是1.0，需要添加alpha通道
	if color.Alpha < 1.0 {
		alpha := int(color.Alpha * 255)
		// 确保hex格式正确（去除#）
		hexValue := strings.TrimPrefix(hex, "#")
		return fmt.Sprintf("#%02X%s", alpha, hexValue)
	}
	
	return hex
}

// sortColorNames 排序颜色名称
func (g *AndroidGenerator) sortColorNames(colors map[string]string) []string {
	names := make([]string, 0, len(colors))
	for name := range colors {
		names = append(names, name)
	}
	
	// 简单的冒泡排序（因为颜色数量通常不多）
	for i := 0; i < len(names)-1; i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	
	return names
}