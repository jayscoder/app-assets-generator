package color

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// IOSGenerator iOS颜色资源生成器
type IOSGenerator struct {
	outputPath string
}

// NewIOSGenerator 创建iOS生成器
func NewIOSGenerator(outputPath string) *IOSGenerator {
	return &IOSGenerator{
		outputPath: outputPath,
	}
}

// iOSColorSet iOS颜色集结构
type iOSColorSet struct {
	Colors []iOSColor `json:"colors"`
	Info   iOSInfo    `json:"info"`
}

// iOSColor iOS颜色定义
type iOSColor struct {
	Color       *iOSColorValue       `json:"color,omitempty"`
	Appearances []iOSAppearance      `json:"appearances,omitempty"`
	Idiom       string               `json:"idiom"`
}

// iOSColorValue iOS颜色值
type iOSColorValue struct {
	ColorSpace string          `json:"color-space"`
	Components iOSComponents   `json:"components"`
}

// iOSComponents iOS颜色组件
type iOSComponents struct {
	Alpha string `json:"alpha"`
	Blue  string `json:"blue"`
	Green string `json:"green"`
	Red   string `json:"red"`
}

// iOSAppearance iOS外观定义
type iOSAppearance struct {
	Appearance string `json:"appearance"`
	Value      string `json:"value"`
}

// iOSInfo iOS信息
type iOSInfo struct {
	Author  string `json:"author"`
	Version int    `json:"version"`
}

// Generate 生成iOS颜色资源
func (g *IOSGenerator) Generate(colors map[string]*ColorDefinition) error {
	// 创建输出目录
	if err := os.MkdirAll(g.outputPath, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	
	// 为每个颜色生成colorset（跳过渐变色）
	for name, color := range colors {
		if color.IsGradient() {
			continue // 跳过渐变色
		}
		if err := g.generateColorSet(g.outputPath, name, color); err != nil {
			return fmt.Errorf("生成颜色 %s 失败: %w", name, err)
		}
	}
	
	return nil
}

// generateColorSet 生成单个颜色集
func (g *IOSGenerator) generateColorSet(outputPath, name string, color *ColorDefinition) error {
	// 创建colorset目录
	colorsetPath := filepath.Join(outputPath, name+".colorset")
	if err := os.MkdirAll(colorsetPath, 0755); err != nil {
		return fmt.Errorf("创建colorset目录失败: %w", err)
	}
	
	// 构建颜色集数据
	colorSet := g.buildColorSet(color)
	
	// 生成Contents.json
	contentsPath := filepath.Join(colorsetPath, "Contents.json")
	file, err := os.Create(contentsPath)
	if err != nil {
		return fmt.Errorf("创建Contents.json失败: %w", err)
	}
	defer file.Close()
	
	// 写入JSON（格式化输出）
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(colorSet); err != nil {
		return fmt.Errorf("写入JSON失败: %w", err)
	}
	
	return nil
}

// buildColorSet 构建iOS颜色集数据
func (g *IOSGenerator) buildColorSet(color *ColorDefinition) iOSColorSet {
	colorSet := iOSColorSet{
		Colors: []iOSColor{},
		Info: iOSInfo{
			Author:  "xcode",
			Version: 1,
		},
	}
	
	// 获取各主题颜色
	defaultColor := color.GetDefault()
	lightColor := color.GetLight()
	darkColor := color.GetDark()
	
	// 添加默认颜色（Any Appearance）
	if defaultColor.Hex != "" {
		colorSet.Colors = append(colorSet.Colors, iOSColor{
			Color: g.buildColorValue(defaultColor),
			Idiom: "universal",
		})
	}
	
	// 如果有不同的light颜色，添加Light Appearance
	if lightColor.Hex != "" && (lightColor.Hex != defaultColor.Hex || lightColor.Alpha != defaultColor.Alpha) {
		colorSet.Colors = append(colorSet.Colors, iOSColor{
			Appearances: []iOSAppearance{
				{
					Appearance: "luminosity",
					Value:      "light",
				},
			},
			Color: g.buildColorValue(lightColor),
			Idiom: "universal",
		})
	}
	
	// 如果有不同的dark颜色，添加Dark Appearance
	if darkColor.Hex != "" && (darkColor.Hex != defaultColor.Hex || darkColor.Alpha != defaultColor.Alpha) {
		colorSet.Colors = append(colorSet.Colors, iOSColor{
			Appearances: []iOSAppearance{
				{
					Appearance: "luminosity",
					Value:      "dark",
				},
			},
			Color: g.buildColorValue(darkColor),
			Idiom: "universal",
		})
	}
	
	return colorSet
}

// buildColorValue 构建iOS颜色值
func (g *IOSGenerator) buildColorValue(color ColorValue) *iOSColorValue {
	r, green, b, _ := hexToRGB(color.Hex)
	
	return &iOSColorValue{
		ColorSpace: "srgb",
		Components: iOSComponents{
			Alpha: formatFloat(color.Alpha),
			Red:   formatFloat(r),
			Green: formatFloat(green),
			Blue:  formatFloat(b),
		},
	}
}