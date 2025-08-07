#!/bin/bash

# build.sh - 构建 App Assets Generator
# 用于编译不同平台的二进制文件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目信息
APP_NAME="app-assets-generator"
VERSION="1.0.0"
BUILD_DIR="build"

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

# 检查Go环境
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go未安装，请先安装Go环境"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    print_info "检测到Go版本: $GO_VERSION"
}

# 清理构建目录
clean() {
    print_info "清理构建目录..."
    rm -rf $BUILD_DIR
    mkdir -p $BUILD_DIR
}

# 下载依赖
download_deps() {
    print_info "下载依赖包..."
    go mod download
    go mod tidy
}

# 构建单个平台
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3
    
    print_info "构建 $GOOS/$GOARCH..."
    
    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w -X 'app-assets-generator/cmd.version=$VERSION'" \
        -o "$OUTPUT" \
        main.go
    
    if [ $? -eq 0 ]; then
        print_info "✓ 构建成功: $OUTPUT"
        # 计算文件大小
        SIZE=$(du -h "$OUTPUT" | cut -f1)
        print_info "  文件大小: $SIZE"
    else
        print_error "✗ 构建失败: $GOOS/$GOARCH"
        return 1
    fi
}

# 构建所有平台
build_all() {
    print_info "开始构建所有平台..."
    
    # macOS
    build_platform "darwin" "amd64" "$BUILD_DIR/${APP_NAME}-darwin-amd64"
    build_platform "darwin" "arm64" "$BUILD_DIR/${APP_NAME}-darwin-arm64"
    
    # Linux
    build_platform "linux" "amd64" "$BUILD_DIR/${APP_NAME}-linux-amd64"
    build_platform "linux" "arm64" "$BUILD_DIR/${APP_NAME}-linux-arm64"
    
    # Windows
    build_platform "windows" "amd64" "$BUILD_DIR/${APP_NAME}-windows-amd64.exe"
    build_platform "windows" "arm64" "$BUILD_DIR/${APP_NAME}-windows-arm64.exe"
}

# 构建当前平台
build_current() {
    print_info "构建当前平台..."
    
    go build \
        -ldflags="-s -w -X 'app-assets-generator/cmd.version=$VERSION'" \
        -o "$APP_NAME" \
        main.go
    
    if [ $? -eq 0 ]; then
        print_info "✓ 构建成功: $APP_NAME"
        SIZE=$(du -h "$APP_NAME" | cut -f1)
        print_info "  文件大小: $SIZE"
    else
        print_error "✗ 构建失败"
        exit 1
    fi
}

# 创建发布包
create_release() {
    print_info "创建发布包..."
    
    # 创建临时目录
    RELEASE_DIR="$BUILD_DIR/release"
    mkdir -p $RELEASE_DIR
    
    # 复制文档
    cp README.md $RELEASE_DIR/
    cp LICENSE $RELEASE_DIR/
    
    # 创建示例配置
    cp colors.yaml $RELEASE_DIR/colors-example.yaml
    
    # 为每个平台创建tar.gz包
    for file in $BUILD_DIR/${APP_NAME}-*; do
        if [ -f "$file" ]; then
            BASENAME=$(basename "$file")
            PLATFORM_DIR="$RELEASE_DIR/$BASENAME"
            mkdir -p "$PLATFORM_DIR"
            
            cp "$file" "$PLATFORM_DIR/"
            cp README.md "$PLATFORM_DIR/"
            cp LICENSE "$PLATFORM_DIR/"
            cp colors.yaml "$PLATFORM_DIR/colors-example.yaml"
            
            # 创建压缩包
            TAR_NAME="${BASENAME}.tar.gz"
            print_info "创建压缩包: $TAR_NAME"
            cd $BUILD_DIR
            tar -czf "$TAR_NAME" -C release "$BASENAME"
            cd - > /dev/null
            
            # 清理临时目录
            rm -rf "$PLATFORM_DIR"
        fi
    done
    
    # 清理release目录
    rm -rf $RELEASE_DIR
    
    print_info "✓ 发布包创建完成"
}

# 显示帮助
show_help() {
    echo "App Assets Generator 构建脚本"
    echo ""
    echo "用法: ./build.sh [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help      显示帮助信息"
    echo "  -c, --current   只构建当前平台"
    echo "  -a, --all       构建所有平台"
    echo "  -r, --release   创建发布包"
    echo "  -v, --version   设置版本号"
    echo "  --clean         只清理构建目录"
    echo ""
    echo "示例:"
    echo "  ./build.sh              # 构建当前平台"
    echo "  ./build.sh -a           # 构建所有平台"
    echo "  ./build.sh -a -r        # 构建所有平台并创建发布包"
    echo "  ./build.sh -v 1.2.0     # 设置版本号为1.2.0并构建"
}

# 主函数
main() {
    BUILD_ALL=false
    CREATE_RELEASE=false
    CLEAN_ONLY=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--current)
                BUILD_ALL=false
                shift
                ;;
            -a|--all)
                BUILD_ALL=true
                shift
                ;;
            -r|--release)
                CREATE_RELEASE=true
                shift
                ;;
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            --clean)
                CLEAN_ONLY=true
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
    print_info "App Assets Generator 构建脚本"
    print_info "版本: $VERSION"
    print_info "========================================="
    
    # 检查环境
    check_go
    
    # 清理构建目录
    clean
    
    if [ "$CLEAN_ONLY" = true ]; then
        print_info "✓ 清理完成"
        exit 0
    fi
    
    # 下载依赖
    download_deps
    
    # 构建
    if [ "$BUILD_ALL" = true ]; then
        build_all
        if [ "$CREATE_RELEASE" = true ]; then
            create_release
        fi
    else
        build_current
    fi
    
    print_info "========================================="
    print_info "✓ 构建完成！"
    if [ "$BUILD_ALL" = true ]; then
        print_info "构建产物位于: $BUILD_DIR/"
        ls -lh $BUILD_DIR/${APP_NAME}-* 2>/dev/null | awk '{print "  " $9 " (" $5 ")"}'
    else
        print_info "可执行文件: ./$APP_NAME"
    fi
    print_info "========================================="
}

# 执行主函数
main "$@"