#!/bin/bash
#
# go-tiny-claw 一键部署脚本
# 适用于 Linux 服务器（Ubuntu/Debian/CentOS）
#
# 使用方法：
#   chmod +x deploy.sh
#   ./deploy.sh [环境]
#   环境可选：dev（开发模式）、prod（生产模式）
#
# 示例：
#   ./deploy.sh prod    # 生产模式部署
#   ./deploy.sh dev     # 开发模式部署

set -e

# 配置参数
APP_NAME="go-tiny-claw"
APP_DIR="/opt/${APP_NAME}"
GIT_REPO="https://github.com/skytodmoon/go-tiny-claw.git"
GO_VERSION="1.22"
SERVICE_NAME="${APP_NAME}.service"
ENV_FILE="/etc/${APP_NAME}.env"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印日志
log() {
    echo -e "${YELLOW}[*]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    exit 1
}

# 检查 root 权限
check_root() {
    if [ "$(id -u)" != "0" ]; then
        log_error "请使用 root 用户或 sudo 运行此脚本"
    fi
}

# 检查并安装依赖
install_dependencies() {
    log "检查系统依赖..."
    
    # 检测包管理器
    if command -v apt-get &>/dev/null; then
        PACKAGE_MANAGER="apt-get"
        UPDATE_CMD="apt-get update -y"
        INSTALL_CMD="apt-get install -y"
    elif command -v yum &>/dev/null; then
        PACKAGE_MANAGER="yum"
        UPDATE_CMD="yum update -y"
        INSTALL_CMD="yum install -y"
    else
        log_error "不支持的包管理器"
    fi
    
    # 更新包列表
    $UPDATE_CMD
    
    # 安装必要依赖
    $INSTALL_CMD git wget curl
    
    log_success "系统依赖安装完成"
}

# 检查并安装 Go
install_go() {
    log "检查 Go 环境..."
    
    if command -v go &>/dev/null; then
        CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        log "当前 Go 版本: ${CURRENT_VERSION}"
        
        # 比较版本
        if [ "$(printf '%s\n' "$GO_VERSION" "$CURRENT_VERSION" | sort -V | head -n1)" = "$GO_VERSION" ]; then
            log_success "Go 版本满足要求"
            return
        fi
    fi
    
    log "开始安装 Go ${GO_VERSION}..."
    
    # 下载并安装 Go
    GO_URL="https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"
    wget -q "${GO_URL}" -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    
    # 设置环境变量
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    echo 'export GOPATH=$HOME/go' >> /etc/profile
    echo 'export PATH=$PATH:$GOPATH/bin' >> /etc/profile
    source /etc/profile
    
    # 验证安装
    if ! command -v go &>/dev/null; then
        log_error "Go 安装失败"
    fi
    
    log_success "Go ${GO_VERSION} 安装完成"
}

# 克隆或更新代码
clone_code() {
    log "获取项目代码..."
    
    if [ -d "${APP_DIR}" ]; then
        log "项目已存在，执行更新..."
        cd "${APP_DIR}"
        git pull origin main
    else
        log "克隆项目..."
        git clone "${GIT_REPO}" "${APP_DIR}"
        cd "${APP_DIR}"
    fi
    
    log_success "代码获取完成"
}

# 编译项目
build_project() {
    log "编译项目..."
    
    cd "${APP_DIR}"
    
    # 设置国内代理加速
    export GOPROXY=https://goproxy.cn,direct
    
    # 下载依赖
    go mod tidy
    
    # 编译
    go build -o "${APP_NAME}" ./cmd/claw
    
    if [ ! -f "${APP_NAME}" ]; then
        log_error "编译失败"
    fi
    
    log_success "项目编译完成"
}

# 创建环境配置文件
create_env_file() {
    log "创建环境配置文件..."
    
    cat > "${ENV_FILE}" << EOF
# go-tiny-claw 环境配置

# 飞书配置


| 变量名 | 说明 | 示例值 |
|--------|------|--------|
| FEISHU_APP_ID | 飞书应用 ID | `xxxx` |
| FEISHU_APP_SECRET | 飞书应用密钥 | `xxxx` |
| FEISHU_VERIFY_TOKEN | 事件验证令牌 | `xxxx` |
| FEISHU_ENCRYPT_KEY | 事件加密密钥 | `xxxx` |
| SILICONFLOW_API_KEY | SiliconFlow API Key | 从官网获取 |

# AI API 配置
SILICONFLOW_API_KEY="your-api-key-here"

# 应用配置
WORK_DIR="/opt/${APP_NAME}"
EOF
    
    log_success "环境配置文件创建完成"
    log "请编辑 ${ENV_FILE} 配置正确的 API Key"
}

# 创建 systemd 服务
create_systemd_service() {
    log "创建 systemd 服务..."
    
    cat > "/etc/systemd/system/${SERVICE_NAME}" << EOF
[Unit]
Description=go-tiny-claw Feishu Bot Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${APP_DIR}
EnvironmentFile=${ENV_FILE}
ExecStart=${APP_DIR}/${APP_NAME}
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
    
    # 重载 systemd
    systemctl daemon-reload
    
    log_success "systemd 服务创建完成"
}

# 启动服务
start_service() {
    log "启动服务..."
    
    systemctl enable "${SERVICE_NAME}"
    systemctl start "${SERVICE_NAME}"
    
    # 检查服务状态
    sleep 2
    if systemctl is-active --quiet "${SERVICE_NAME}"; then
        log_success "服务启动成功"
        log "服务状态:"
        systemctl status "${SERVICE_NAME}"
    else
        log_error "服务启动失败"
    fi
}

# 开发模式部署
deploy_dev() {
    log "开始开发模式部署..."
    
    install_dependencies
    install_go
    clone_code
    build_project
    create_env_file
    
    log_success "开发模式部署完成"
    log "启动方式: cd ${APP_DIR} && source ${ENV_FILE} && ./${APP_NAME}"
}

# 生产模式部署
deploy_prod() {
    log "开始生产模式部署..."
    
    install_dependencies
    install_go
    clone_code
    build_project
    create_env_file
    create_systemd_service
    start_service
    
    log_success "生产模式部署完成"
    log "服务已启动，监听端口: 48080"
    log "配置文件: ${ENV_FILE}"
    log "日志查看: journalctl -u ${SERVICE_NAME} -f"
}

# 主函数
main() {
    check_root
    
    ENV_MODE="${1:-prod}"
    
    case "${ENV_MODE}" in
        dev)
            deploy_dev
            ;;
        prod)
            deploy_prod
            ;;
        *)
            log_error "无效的环境参数: ${ENV_MODE}"
            log "使用方法: $0 [dev|prod]"
            ;;
    esac
}

main "$@"
