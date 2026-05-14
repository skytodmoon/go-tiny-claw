# go-tiny-claw 部署指南

本文档介绍如何在 Linux 服务器上部署 go-tiny-claw 飞书机器人服务。

## 目录

1. [环境要求](#环境要求)
2. [编译项目](#编译项目)
3. [创建 systemd 服务](#创建-systemd-服务)
4. [启动与管理服务](#启动与管理服务)
5. [配置飞书](#配置飞书)
6. [常见问题](#常见问题)

---

## 环境要求

| 依赖 | 版本要求 |
|------|----------|
| Go | 1.20+ |
| Linux | Ubuntu/Debian/CentOS |
| 网络 | 可访问飞书 API 和 SiliconFlow |

---

## 编译项目

### 1. 安装 Go 环境

```bash
# Ubuntu/Debian
sudo apt-get update && sudo apt-get install -y golang

# CentOS/RHEL
sudo yum install -y golang

# 验证安装
go version
```

### 2. 设置国内代理（可选但推荐）

```bash
# 临时生效
export GOPROXY=https://goproxy.cn,direct

# 永久生效（写入 ~/.bashrc）
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
source ~/.bashrc
```

### 3. 克隆项目

```bash
cd /home/admin
git clone https://github.com/skytodmoon/go-tiny-claw.git
cd go-tiny-claw
```

### 4. 编译项目

```bash
# 下载依赖
go mod tidy

# 编译
go build -o go-tiny-claw ./cmd/claw

# 验证编译结果
ls -la go-tiny-claw
```

---

## 创建 systemd 服务

### 创建服务配置文件

```bash
sudo cat > /etc/systemd/system/go-tiny-claw.service << 'EOF'
[Unit]
Description=go-tiny-claw Feishu Bot Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/home/admin/go-tiny-claw
Environment="FEISHU_APP_ID=xxxx"
Environment="FEISHU_APP_SECRET=xxxx"
Environment="FEISHU_VERIFY_TOKEN=xxxx"
Environment="FEISHU_ENCRYPT_KEY=xxxx"
Environment="SILICONFLOW_API_KEY=xxxx"
ExecStart=/home/admin/go-tiny-claw/go-tiny-claw
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

### 设置文件权限

```bash
sudo chmod 644 /etc/systemd/system/go-tiny-claw.service
```

---

## 启动与管理服务

### 启动服务

```bash
# 重载 systemd 配置
sudo systemctl daemon-reload

# 设置开机自启
sudo systemctl enable go-tiny-claw

# 启动服务
sudo systemctl start go-tiny-claw

# 查看状态
sudo systemctl status go-tiny-claw
```

### 常用管理命令

| 命令 | 说明 |
|------|------|
| `systemctl start go-tiny-claw` | 启动服务 |
| `systemctl stop go-tiny-claw` | 停止服务 |
| `systemctl restart go-tiny-claw` | 重启服务 |
| `systemctl status go-tiny-claw` | 查看状态 |
| `systemctl enable go-tiny-claw` | 设置开机自启 |
| `systemctl disable go-tiny-claw` | 取消开机自启 |
| `journalctl -u go-tiny-claw -f` | 查看实时日志 |

---

## 配置飞书

### 在飞书开发者后台配置

1. **进入飞书开发者平台**：https://open.feishu.cn/app
2. **找到你的应用**（App ID: `xxxx`）
3. **配置事件订阅**：
   - 请求 URL: `http://your-server-ip:48080/webhook/event`
   - 加密密钥: `xxxx`
   - 验证令牌: `xxxx`
4. **添加权限**：
   - `im:message:send` - 发送消息

---

## 常见问题

### 1. Failed to enable unit: Invalid argument

**原因**：服务文件语法错误或权限问题

**解决方法**：

```bash
# 验证服务文件语法
sudo systemd-analyze verify go-tiny-claw.service

# 重新创建服务文件（使用 cat 方式避免复制问题）
sudo rm /etc/systemd/system/go-tiny-claw.service
sudo cat > /etc/systemd/system/go-tiny-claw.service << 'EOF'
[Unit]
Description=go-tiny-claw Feishu Bot Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/home/admin/go-tiny-claw
Environment="FEISHU_APP_ID=xxxx"
Environment="FEISHU_APP_SECRET=xxxx"
Environment="FEISHU_VERIFY_TOKEN=xxxx"
Environment="FEISHU_ENCRYPT_KEY=xxxx"
Environment="SILICONFLOW_API_KEY=xxxx"
ExecStart=/home/admin/go-tiny-claw/go-tiny-claw
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo chmod 644 /etc/systemd/system/go-tiny-claw.service
sudo systemctl daemon-reload
```

### 2. 服务启动失败

**检查日志**：

```bash
journalctl -u go-tiny-claw -f
```

**常见原因**：
- SILICONFLOW_API_KEY 未配置
- 端口 48080 被占用
- 飞书凭证错误

### 3. 飞书消息发送失败

**检查飞书配置**：

```bash
# 测试飞书 API
curl -X POST https://open.feishu.cn/open-apis/auth/v3/app_access_token \
  -H "Content-Type: application/json" \
  -d '{"app_id":" xxxx","app_secret":"xxxx"}'
```

**确保应用已发布**：在飞书开发者后台确认应用状态为"已发布"

---

## 配置参考

### 环境变量列表

| 变量名 | 说明 | 示例值 |
|--------|------|--------|
| FEISHU_APP_ID | 飞书应用 ID | `xxxx` |
| FEISHU_APP_SECRET | 飞书应用密钥 | `xxxx` |
| FEISHU_VERIFY_TOKEN | 事件验证令牌 | `xxxx` |
| FEISHU_ENCRYPT_KEY | 事件加密密钥 | `xxxx` |
| SILICONFLOW_API_KEY | SiliconFlow API Key | 从官网获取 |

### 服务监听端口

默认端口：`48080`

---

## 更新项目

```bash
# 停止服务
sudo systemctl stop go-tiny-claw

# 拉取最新代码
cd /home/admin/go-tiny-claw
git pull

# 重新编译
go build -o go-tiny-claw ./cmd/claw

# 启动服务
sudo systemctl start go-tiny-claw
```

---

## 卸载服务

```bash
# 停止服务
sudo systemctl stop go-tiny-claw

# 取消开机自启
sudo systemctl disable go-tiny-claw

# 删除服务文件
sudo rm /etc/systemd/system/go-tiny-claw.service

# 重载 systemd
sudo systemctl daemon-reload

# 删除项目文件（可选）
rm -rf /home/admin/go-tiny-claw
```
