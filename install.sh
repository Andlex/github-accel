#!/bin/bash
set -e
echo "=== GitHub DNS 加速 ==="
sudo mkdir -p /etc/systemd/resolved.conf.d/
sudo tee /etc/systemd/resolved.conf.d/github.conf > /dev/null << 'CONF'
[Resolve]
DNS=223.5.5.5#dns.alidns.com
DNSOverTLS=yes
Domains=~github.com
CONF
sudo systemctl restart systemd-resolved
echo "配置完成"
resolvectl status | grep -E "DNSOverTLS|DNS Servers"
