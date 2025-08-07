package color

import (
	"fmt"
	"os"
	
	"gopkg.in/yaml.v3"
)

// ParseYAML 解析YAML颜色配置文件
func ParseYAML(filePath string) (map[string]*ColorDefinition, error) {
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	
	// 解析YAML
	var colors map[string]*ColorDefinition
	if err := yaml.Unmarshal(data, &colors); err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}
	
	// 验证颜色值
	for name, color := range colors {
		if err := validateColor(name, color); err != nil {
			return nil, err
		}
	}
	
	return colors, nil
}

// validateColor 验证颜色定义
func validateColor(name string, color *ColorDefinition) error {
	if color == nil {
		return fmt.Errorf("颜色 %s 定义为空", name)
	}
	
	// 渐变色暂时跳过
	if color.IsGradient() {
		return nil
	}
	
	// 简单颜色验证
	if color.IsSimple() {
		if !isValidHex(color.Hex) {
			return fmt.Errorf("颜色 %s 的hex值无效: %s", name, color.Hex)
		}
		if color.Alpha < 0 || color.Alpha > 1 {
			return fmt.Errorf("颜色 %s 的alpha值必须在0-1之间: %f", name, color.Alpha)
		}
		return nil
	}
	
	// 主题颜色验证
	if color.Default == nil && color.Light == nil && color.Dark == nil {
		return fmt.Errorf("颜色 %s 必须至少定义一个主题颜色", name)
	}
	
	// 验证各主题颜色
	if color.Default != nil {
		if !isValidHex(color.Default.Hex) {
			return fmt.Errorf("颜色 %s 的default.hex值无效: %s", name, color.Default.Hex)
		}
		if color.Default.Alpha < 0 || color.Default.Alpha > 1 {
			return fmt.Errorf("颜色 %s 的default.alpha值必须在0-1之间: %f", name, color.Default.Alpha)
		}
	}
	
	if color.Light != nil {
		if !isValidHex(color.Light.Hex) {
			return fmt.Errorf("颜色 %s 的light.hex值无效: %s", name, color.Light.Hex)
		}
		if color.Light.Alpha < 0 || color.Light.Alpha > 1 {
			return fmt.Errorf("颜色 %s 的light.alpha值必须在0-1之间: %f", name, color.Light.Alpha)
		}
	}
	
	if color.Dark != nil {
		if !isValidHex(color.Dark.Hex) {
			return fmt.Errorf("颜色 %s 的dark.hex值无效: %s", name, color.Dark.Hex)
		}
		if color.Dark.Alpha < 0 || color.Dark.Alpha > 1 {
			return fmt.Errorf("颜色 %s 的dark.alpha值必须在0-1之间: %f", name, color.Dark.Alpha)
		}
	}
	
	return nil
}

// isValidHex 验证十六进制颜色值
func isValidHex(hex string) bool {
	if len(hex) != 7 || hex[0] != '#' {
		return false
	}
	
	for i := 1; i < 7; i++ {
		c := hex[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	
	return true
}