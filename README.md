# APP Assets Generator

一个用于生成Android和iOS应用资源文件的命令行工具，支持颜色和图片资源的批量生成。

## 功能特性

- 🎨 **颜色资源生成** - 从YAML配置文件批量生成iOS和Android的颜色资源
- 🖼️ **图片资源生成** - 自动处理@2x、@3x等多分辨率图片资源
- 🌓 **深色模式支持** - 支持Light/Dark主题的颜色配置
- 📱 **多平台支持** - 同时支持iOS和Android平台
- ⚡ **批量处理** - 支持批量处理多个资源文件

## 安装

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/yourusername/app-assets-generator.git
cd app-assets-generator

# 安装依赖
go mod download

# 构建
go build -o app-assets-generator main.go
```

### 使用Go Install

```bash
go install github.com/yourusername/app-assets-generator@latest
```

## 使用方式

### 基本命令

```bash
# 查看帮助
app-assets-generator --help

# 查看版本
app-assets-generator --version

# 生成颜色资源
app-assets-generator color --input colors.yaml --output output/colors --platform ios
app-assets-generator color --input colors.yaml --output output/colors --platform android

# 生成图片资源
app-assets-generator image --input icons/ --output output/images --platform ios
app-assets-generator image --input icons/ --output output/images --platform android
```

### 生成颜色资源

从YAML配置文件生成平台特定的颜色资源：

```bash
# iOS平台
app-assets-generator color --input=colors.yaml --output=output/colors-ios --platform=ios

# Android平台
app-assets-generator color --input=colors.yaml --output=output/colors-android --platform=android

# 同时生成两个平台
app-assets-generator color --input=colors.yaml --output=output/ --platform=all
```

#### 颜色配置文件格式

`colors.yaml` 示例：

```yaml
# 基础颜色定义
color_primary:
  default:
    hex: "#34a3f4"
    alpha: 1.0
  light:
    hex: "#34a3f4"
    alpha: 1.0
  dark:
    hex: "#5db6f6"
    alpha: 1.0

# 简单颜色定义（不区分主题）
color_simple:
  hex: "#ff0000"
  alpha: 1.0
```

#### iOS输出格式

生成的iOS颜色资源直接位于指定的输出目录：
- `[color-name].colorset/Contents.json`
```json
{
  "colors" : [
    {
      "color" : {
        "color-space" : "srgb",
        "components" : {
          "alpha" : "0.000",
          "blue" : "0.000",
          "green" : "0.000",
          "red" : "0.000"
        }
      },
      "idiom" : "universal"
    },
    {
      "appearances" : [
        {
          "appearance" : "luminosity",
          "value" : "light"
        }
      ],
      "color" : {
        "color-space" : "srgb",
        "components" : {
          "alpha" : "0.000",
          "blue" : "0.000",
          "green" : "0.000",
          "red" : "0.000"
        }
      },
      "idiom" : "universal"
    },
    {
      "appearances" : [
        {
          "appearance" : "luminosity",
          "value" : "dark"
        }
      ],
      "color" : {
        "color-space" : "srgb",
        "components" : {
          "alpha" : "0.000",
          "blue" : "0.000",
          "green" : "0.000",
          "red" : "0.000"
        }
      },
      "idiom" : "universal"
    }
  ],
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}
```

#### Android输出格式

生成的Android颜色资源：
- `values/colors.xml` - 默认颜色
- `values-night/colors.xml` - 深色模式颜色

```xml
<!-- values/colors.xml -->
<resources>
    <color name="color_primary">#34a3f4</color>
</resources>

<!-- values-night/colors.xml -->
<resources>
    <color name="color_primary">#5db6f6</color>
</resources>
```

### 生成图片资源

自动处理多分辨率图片并生成平台特定的资源：

```bash
# iOS平台
app-assets-generator image --input=icons/ --output=output/images-ios --platform=ios

# Android平台
app-assets-generator image --input=icons/ --output=output/images-android --platform=android

# 批量处理
app-assets-generator image --input=icons/ --output=output/ --platform=all
```

#### iOS图片资源

生成的资源直接位于指定的输出目录：
- `[image-name].imageset/Contents.json`
- 自动识别 @2x、@3x 后缀的图片文件

Contents.json 示例：
```json
{
  "images" : [
    {
      "filename" : "icon.png",
      "idiom" : "universal",
      "scale" : "1x"
    },
    {
      "filename" : "icon@2x.png",
      "idiom" : "universal",
      "scale" : "2x"
    },
    {
      "filename" : "icon@3x.png",
      "idiom" : "universal",
      "scale" : "3x"
    }
  ],
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}
```

#### Android图片资源

生成的资源分布在不同的drawable目录：
- `drawable-mdpi/` - 1x 图片
- `drawable-hdpi/` - 1.5x 图片
- `drawable-xhdpi/` - 2x 图片
- `drawable-xxhdpi/` - 3x 图片
- `drawable-xxxhdpi/` - 4x 图片

## 配置文件

### 全局配置 (.app-assets-generator.yaml)

可以在项目根目录创建配置文件：

```yaml
# 默认输出目录
output_dir: ./output

# 默认平台
platform: all

# iOS特定配置
ios:
  deployment_target: "13.0"

# Android特定配置
android:
  res_path: src/main/res
  min_sdk_version: 21
```

## 开发

### 项目结构

```
app-assets-generator/
├── main.go              # 主程序入口
├── cmd/                 # 命令行处理
│   ├── root.go         # 根命令
│   ├── color.go        # 颜色生成命令
│   └── image.go        # 图片生成命令
├── pkg/                 # 核心功能
│   ├── color/          # 颜色处理
│   │   ├── parser.go   # YAML解析
│   │   ├── ios.go      # iOS颜色生成
│   │   └── android.go  # Android颜色生成
│   ├── image/          # 图片处理
│   │   ├── scanner.go  # 图片扫描
│   │   ├── ios.go      # iOS图片生成
│   │   └── android.go  # Android图片生成
│   └── utils/          # 工具函数
├── colors.yaml         # 颜色配置示例
└── icons/              # 图标资源示例
```

### 贡献指南

欢迎提交 Pull Request 或创建 Issue！

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 致谢

- 感谢所有贡献者
- 基于 [Cobra](https://github.com/spf13/cobra) 构建命令行界面
- 使用 [YAML](https://github.com/go-yaml/yaml) 进行配置解析