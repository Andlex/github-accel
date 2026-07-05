#!/bin/bash
set -e
echo "=== GitHub Accelerator 安装 ==="

# 安装二进制
install -m 755 github-accel /usr/local/bin/github-accel
echo "二进制已安装到 /usr/local/bin/github-accel"

# 配置 hosts
bash setup.sh

# 创建 polkit policy
sudo tee /usr/share/polkit-1/actions/com.github-accel.policy > /dev/null << 'POLICY'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE policyconfig PUBLIC
 "-//freedesktop//DTD PolicyKit Policy Configuration 1.0//EN"
 "http://www.freedesktop.org/standards/PolicyKit/1/policyconfig.dtd">
<policyconfig>
  <action id="com.github-accel.run">
    <description>Run GitHub Accelerator</description>
    <message>Authentication is required to run GitHub Accelerator</message>
    <defaults>
      <allow_any>auth_admin</allow_any>
      <allow_inactive>auth_admin</allow_inactive>
      <allow_active>auth_admin</allow_active>
    </defaults>
    <annotate key="org.freedesktop.policykit.exec.path">/usr/local/bin/github-accel</annotate>
    <annotate key="org.freedesktop.policykit.exec.allow_gui">true</annotate>
  </action>
</policyconfig>
POLICY
echo "Polkit 策略已创建"

# 创建桌面快捷方式
mkdir -p ~/Desktop
cat > ~/Desktop/GitHub加速器.desktop << 'DESKTOP'
[Desktop Entry]
Name=GitHub 加速器
Comment=轻量 GitHub 加速器
Exec=bash -c 'pkexec /usr/local/bin/github-accel'
Icon=network-transmit-receive
Terminal=true
Type=Application
Categories=Network;
DESKTOP
chmod +x ~/Desktop/GitHub加速器.desktop
echo "桌面快捷方式已创建"

echo ""
echo "=== 安装完成 ==="
echo "双击桌面图标或运行: sudo github-accel"
