# github-accel

极简 GitHub 加速器。Go 静态二进制，零运行时依赖。

## 特性

- **零依赖**：静态链接，无需 Python/Node/.NET
- **极小体积**：2MB 二进制
- **极低内存**：~5MB 运行时（vs Watt Toolkit 665MB）
- **跨平台**：Linux/macOS/Windows

## 使用

```bash
# Linux (需 root)
sudo ./github-accel

# macOS
sudo ./github-accel

# Windows (管理员)
github-accel.exe
```

## 构建

```bash
# 本地构建
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel .

# 交叉编译
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel-linux .
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel-mac .
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel.exe .
```

## 文件

- `main.go` — 源码
- `github-accel` — 编译后的二进制
- `setup.sh` — hosts 配置脚本
- `accel.service` — systemd 服务文件

## 原理

1. `/etc/hosts` 将 GitHub 域名指向 `127.0.0.1`
2. 本地代理监听 443 端口
3. 解析 TLS SNI 域名
4. 转发到真实 GitHub 服务器

## License

MIT
