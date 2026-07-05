# github-accel

轻量 GitHub 加速器。Go 静态二进制，开箱即用。

## 特性

- **零依赖**：静态链接，下载即用
- **体积小**：2MB
- **内存低**：~5MB
- **跨平台**：Linux / macOS / Windows

## 快速开始

```bash
# Linux / macOS
sudo ./github-accel

# Windows (管理员)
github-accel.exe
```

首次运行前，执行 hosts 配置：

```bash
sudo bash setup.sh
```

## 构建

```bash
# 本地
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel .

# 交叉编译
GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel-linux .
GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel-mac .
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel.exe .
```

## 原理

1. `/etc/hosts` 将 GitHub 域名指向 `127.0.0.1`
2. 本地代理监听 443 端口
3. 解析 TLS SNI 域名
4. 转发到真实 GitHub 服务器

## 文件

- `main.go` — 源码
- `setup.sh` — hosts 配置脚本
- `accel.service` — systemd 服务文件

## License

MIT
