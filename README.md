# github-accel

轻量 GitHub 加速器。Go 静态二进制，开箱即用。

## 特性

- **零依赖**：静态链接，下载即用
- **体积小**：2MB
- **内存低**：~2.4MB
- **跨平台**：Linux / macOS / Windows

## 一键安装

```bash
git clone https://github.com/Andlex/github-accel.git
cd github-accel
sudo bash install.sh
```

安装后：
- 自动配置 hosts
- 创建桌面快捷方式
- 双击图标，输入密码即可使用

## 手动使用

```bash
# Linux / macOS
sudo ./github-accel

# Windows
github-accel.exe
```

## 构建

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o github-accel .
```

## 原理

1. `/etc/hosts` 将 GitHub 域名指向 `127.0.0.1`
2. 本地代理监听 443 端口
3. 解析 TLS SNI 域名
4. 转发到真实 GitHub 服务器

## License

MIT
