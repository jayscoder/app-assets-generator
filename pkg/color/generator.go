package color

import (
	"fmt"
	"path/filepath"
)

// Generator 颜色资源生成器
type Generator struct {
	inputPath  string                      // 输入文件路径
	outputPath string                      // 输出目录路径
	colors     map[string]*ColorDefinition // 解析后的颜色数据
}

// NewGenerator 创建新的生成器
func NewGenerator(inputPath, outputPath string) *Generator {
	return &Generator{
		inputPath:  inputPath,
		outputPath: outputPath,
	}
}

// GenerateIOS 生成iOS颜色资源
func (g *Generator) GenerateIOS() error {
	// 解析颜色配置
	if err := g.parseColors(); err != nil {
		return err
	}
	
	// 生成iOS资源
	iosGen := NewIOSGenerator(g.outputPath)
	return iosGen.Generate(g.colors)
}

// GenerateAndroid 生成Android颜色资源
func (g *Generator) GenerateAndroid() error {
	// 解析颜色配置
	if err := g.parseColors(); err != nil {
		return err
	}
	
	// 生成Android资源
	androidGen := NewAndroidGenerator(g.outputPath)
	return androidGen.Generate(g.colors)
}

// parseColors 解析颜色配置（如果还没有解析）
func (g *Generator) parseColors() error {
	if g.colors != nil {
		return nil // 已经解析过了
	}
	
	colors, err := ParseYAML(g.inputPath)
	if err != nil {
		return fmt.Errorf("解析颜色配置失败: %w", err)
	}
	
	g.colors = colors
	return nil
}

// hexToRGB 将十六进制颜色转换为RGB值
func hexToRGB(hex string) (r, g, b float64, err error) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, fmt.Errorf("无效的hex颜色值: %s", hex)
	}
	
	var ri, gi, bi int
	_, err = fmt.Sscanf(hex, "#%02x%02x%02x", &ri, &gi, &bi)
	if err != nil {
		return 0, 0, 0, err
	}
	
	// 转换为0-1的浮点数
	r = float64(ri) / 255.0
	g = float64(gi) / 255.0
	b = float64(bi) / 255.0
	
	return r, g, b, nil
}

// formatFloat 格式化浮点数，去除不必要的小数位
func formatFloat(f float64) string {
	s := fmt.Sprintf("%.3f", f)
	// 去除末尾的0
	for len(s) > 1 && s[len(s)-1] == '0' && s[len(s)-2] != '.' {
		s = s[:len(s)-1]
	}
	return s
}

// ensureDir 确保目录存在
func ensureDir(path string) error {
	dir := filepath.Dir(path)
	return ensureDirPath(dir)
}

// ensureDirPath 确保目录路径存在
func ensureDirPath(dir string) error {
	return nil // 将在具体生成器中实现
}