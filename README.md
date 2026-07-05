# github-accel

GitHub DNS 加速方案。通过 DNS over TLS 防止 DNS 劫持，无需代理进程。

## 原理

ISP 劫持 DNS 返回错误的 GitHub IP → 使用加密 DNS (DoT) 绕过劫持

## 一键配置

```bash
sudo bash install.sh
```

## 手动配置

```bash
sudo mkdir -p /etc/systemd/resolved.conf.d/
sudo tee /etc/systemd/resolved.conf.d/github.conf << 'EOF'
[Resolve]
DNS=223.5.5.5#dns.alidns.com
DNSOverTLS=yes
Domains=~github.com
