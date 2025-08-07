#!/bin/bash

# install.sh - 安装 App Assets Generator
# 自动下载并安装到系统路径

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
APP_NAME="app-assets-generator"
GITHUB_REPO="yourusername/app-assets-generator"  # 需要更新为实际的GitHub仓库
DEFAULT_VERSION="latest"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/${APP_NAME}-install"

# 打印带颜色的信息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检测操作系统和架构
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    # 转换架构名称
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "不支持的架构: $ARCH"
            exit 1
            ;;
    esac
    
    # 转换操作系统名称
    case $OS in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            print_error "请使用Windows安装脚本"
            exit 1
            ;;
        *)
            print_error "不支持的操作系统: $OS"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}-${ARCH}"
    print_info "检测到平台: $PLATFORM"
}

# 检查依赖
check_deps() {
    print_step "检查依赖..."
    
    # 检查curl或wget
    if command -v curl &> /dev/null; then
        DOWNLOADER="curl"
    elif command -v wget &> /dev/null; then
        DOWNLOADER="wget"
    else
        print_error "需要安装 curl 或 wget"
        exit 1
    fi
    
    # 检查tar
    if ! command -v tar &> /dev/null; then
        print_error "需要安装 tar"
        exit 1
    fi
    
    print_info "使用 $DOWNLOADER 下载"
}

# 检查权限
check_permissions() {
    print_step "检查安装权限..."
    
    if [ ! -w "$INSTALL_DIR" ]; then
        print_warning "没有写入权限: $INSTALL_DIR"
        print_info "尝试使用 sudo..."
        NEED_SUDO=true
    else
        NEED_SUDO=false
    fi
}

# 从源码安装
install_from_source() {
    print_step "从源码安装..."
    
    # 检查Go环境
    if ! command -v go &> /dev/null; then
        print_error "未检测到Go环境，无法从源码安装"
        print_info "请先安装Go: https://golang.org/dl/"
        exit 1
    fi
    
    # 检查git
    if ! command -v git &> /dev/null; then
        print_error "未检测到git，无法从源码安装"
        exit 1
    fi
    
    # 创建临时目录
    mkdir -p $TEMP_DIR
    cd $TEMP_DIR
    
    # 克隆仓库
    print_info "克隆仓库..."
    if [ -d "$APP_NAME" ]; then
        rm -rf "$APP_NAME"
    fi
    git clone "https://github.com/${GITHUB_REPO}.git" "$APP_NAME"
    cd "$APP_NAME"
    
    # 构建
    print_info "构建程序..."
    go mod download
    go build -ldflags="-s -w" -o "$APP_NAME" main.go
    
    # 安装
    print_info "安装到 $INSTALL_DIR..."
    if [ "$NEED_SUDO" = true ]; then
        sudo mv "$APP_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$APP_NAME"
    else
        mv "$APP_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$APP_NAME"
    fi
    
    # 清理
    cd /
    rm -rf $TEMP_DIR
}

# 从发布包安装
install_from_release() {
    local VERSION=$1
    print_step "从发布包安装 (版本: $VERSION)..."
    
    # 构建下载URL
    if [ "$VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/${APP_NAME}-${PLATFORM}.tar.gz"
    else
        DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${APP_NAME}-${PLATFORM}.tar.gz"
    fi
    
    # 创建临时目录
    mkdir -p $TEMP_DIR
    cd $TEMP_DIR
    
    # 下载
    print_info "下载 $DOWNLOAD_URL..."
    if [ "$DOWNLOADER" = "curl" ]; then
        curl -L -o "${APP_NAME}.tar.gz" "$DOWNLOAD_URL"
    else
        wget -O "${APP_NAME}.tar.gz" "$DOWNLOAD_URL"
    fi
    
    # 解压
    print_info "解压..."
    tar -xzf "${APP_NAME}.tar.gz"
    
    # 查找可执行文件
    BINARY=$(find . -type f -name "$APP_NAME*" -executable | head -n 1)
    if [ -z "$BINARY" ]; then
        # 如果没找到可执行文件，尝试查找任何二进制文件
        BINARY=$(find . -type f -name "$APP_NAME*" | head -n 1)
    fi
    
    if [ -z "$BINARY" ]; then
        print_error "未找到可执行文件"
        exit 1
    fi
    
    # 安装
    print_info "安装到 $INSTALL_DIR..."
    if [ "$NEED_SUDO" = true ]; then
        sudo cp "$BINARY" "$INSTALL_DIR/$APP_NAME"
        sudo chmod +x "$INSTALL_DIR/$APP_NAME"
    else
        cp "$BINARY" "$INSTALL_DIR/$APP_NAME"
        chmod +x "$INSTALL_DIR/$APP_NAME"
    fi
    
    # 清理
    cd /
    rm -rf $TEMP_DIR
}

# 本地安装（用于开发）
install_local() {
    print_step "本地安装..."
    
    if [ ! -f "main.go" ]; then
        print_error "未找到main.go，请在项目根目录运行"
        exit 1
    fi
    
    # 检查Go环境
    if ! command -v go &> /dev/null; then
        print_error "未检测到Go环境"
        exit 1
    fi
    
    # 构建
    print_info "构建程序..."
    go mod download
    go build -ldflags="-s -w" -o "$APP_NAME" main.go
    
    if [ ! -f "$APP_NAME" ]; then
        print_error "构建失败"
        exit 1
    fi
    
    # 安装
    print_info "安装到 $INSTALL_DIR..."
    if [ "$NEED_SUDO" = true ]; then
        sudo cp "$APP_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$APP_NAME"
    else
        cp "$APP_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$APP_NAME"
    fi
}

# 验证安装
verify_installation() {
    print_step "验证安装..."
    
    # 检查文件是否存在
    if [ ! -f "$INSTALL_DIR/$APP_NAME" ]; then
        print_error "安装失败：文件不存在"
        exit 1
    fi
    
    # 检查是否可执行
    if [ ! -x "$INSTALL_DIR/$APP_NAME" ]; then
        print_error "安装失败：文件不可执行"
        exit 1
    fi
    
    # 检查版本
    if command -v "$APP_NAME" &> /dev/null; then
        VERSION_OUTPUT=$("$APP_NAME" --version 2>&1 || true)
        print_info "已安装版本: $VERSION_OUTPUT"
    else
        print_warning "程序已安装到 $INSTALL_DIR，但不在PATH中"
        print_info "请将 $INSTALL_DIR 添加到PATH环境变量"
    fi
}

# 卸载
uninstall() {
    print_step "卸载 $APP_NAME..."
    
    if [ -f "$INSTALL_DIR/$APP_NAME" ]; then
        if [ "$NEED_SUDO" = true ]; then
            sudo rm -f "$INSTALL_DIR/$APP_NAME"
        else
            rm -f "$INSTALL_DIR/$APP_NAME"
        fi
        print_info "✓ 已卸载"
    else
        print_warning "未找到安装文件"
    fi
}

# 显示帮助
show_help() {
    echo "App Assets Generator 安装脚本"
    echo ""
    echo "用法: ./install.sh [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -v, --version       指定版本号 (默认: latest)"
    echo "  -d, --dir           指定安装目录 (默认: /usr/local/bin)"
    echo "  -s, --source        从源码安装"
    echo "  -l, --local         从本地源码安装（开发用）"
    echo "  -u, --uninstall     卸载"
    echo ""
    echo "示例:"
    echo "  ./install.sh                    # 安装最新版本"
    echo "  ./install.sh -v v1.0.0          # 安装指定版本"
    echo "  ./install.sh -d ~/bin           # 安装到指定目录"
    echo "  ./install.sh -s                 # 从源码安装"
    echo "  ./install.sh -l                 # 从本地构建并安装"
    echo "  ./install.sh -u                 # 卸载"
}

# 主函数
main() {
    VERSION="$DEFAULT_VERSION"
    FROM_SOURCE=false
    FROM_LOCAL=true
    UNINSTALL=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -d|--dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            -s|--source)
                FROM_SOURCE=true
                shift
                ;;
            -l|--local)
                FROM_LOCAL=true
                shift
                ;;
            -u|--uninstall)
                UNINSTALL=true
                shift
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "========================================="
    print_info "App Assets Generator 安装脚本"
    print_info "========================================="
    
    # 检测平台
    detect_platform
    
    # 检查权限
    check_permissions
    
    # 执行操作
    if [ "$UNINSTALL" = true ]; then
        uninstall
    elif [ "$FROM_LOCAL" = true ]; then
        install_local
        verify_installation
    elif [ "$FROM_SOURCE" = true ]; then
        check_deps
        install_from_source
        verify_installation
    else
        check_deps
        install_from_release "$VERSION"
        verify_installation
    fi
    
    if [ "$UNINSTALL" = false ]; then
        print_info "========================================="
        print_info "✓ 安装成功！"
        print_info "运行 '$APP_NAME --help' 查看使用方法"
        print_info "========================================="
    fi
}

# 执行主函数
main "$@"