package color

// ColorValue 颜色值定义
type ColorValue struct {
	Hex   string  `yaml:"hex"`   // 十六进制颜色值
	Alpha float64 `yaml:"alpha"` // 透明度 0.0-1.0
}

// ColorDefinition 颜色定义（支持主题）
type ColorDefinition struct {
	// 简单模式（不区分主题）
	Hex   string  `yaml:"hex,omitempty"`
	Alpha float64 `yaml:"alpha,omitempty"`
	
	// 主题模式
	Default *ColorValue `yaml:"default,omitempty"` // 默认颜色
	Light   *ColorValue `yaml:"light,omitempty"`   // 浅色主题
	Dark    *ColorValue `yaml:"dark,omitempty"`    // 深色主题
	
	// 渐变模式（暂时忽略）
	Type     string                   `yaml:"type,omitempty"`     // 渐变类型
	Angle    string                   `yaml:"angle,omitempty"`    // 渐变角度
	Opacity  float64                  `yaml:"opacity,omitempty"`  // 渐变透明度
	Stops    []map[string]interface{} `yaml:"stops,omitempty"`    // 渐变停止点
}

// IsSimple 判断是否为简单颜色（不区分主题）
func (c *ColorDefinition) IsSimple() bool {
	return c.Hex != "" && c.Default == nil && c.Type == ""
}

// IsGradient 判断是否为渐变色
func (c *ColorDefinition) IsGradient() bool {
	return c.Type != ""
}

// GetDefault 获取默认颜色
func (c *ColorDefinition) GetDefault() ColorValue {
	if c.IsSimple() {
		return ColorValue{Hex: c.Hex, Alpha: c.Alpha}
	}
	if c.Default != nil {
		return *c.Default
	}
	// 如果没有default，尝试使用light
	if c.Light != nil {
		return *c.Light
	}
	return ColorValue{}
}

// GetLight 获取浅色主题颜色
func (c *ColorDefinition) GetLight() ColorValue {
	if c.IsSimple() {
		return ColorValue{Hex: c.Hex, Alpha: c.Alpha}
	}
	if c.Light != nil {
		return *c.Light
	}
	// 如果没有light，使用default
	return c.GetDefault()
}

// GetDark 获取深色主题颜色
func (c *ColorDefinition) GetDark() ColorValue {
	if c.IsSimple() {
		return ColorValue{Hex: c.Hex, Alpha: c.Alpha}
	}
	if c.Dark != nil {
		return *c.Dark
	}
	// 如果没有dark，使用default
	return c.GetDefault()
}