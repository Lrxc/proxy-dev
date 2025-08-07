proxy-dev

系统级抓包, 请求重定向,拦截,修改等

# 功能列表:

- 网络拦截重定向

# 打包

安装依赖
```shell
#Ubuntu/Debian
sudo apt-get install -y libx11-dev libgl1-mesa-dev xorg-dev
sudo apt-get install -y libxxf86vm-dev

#CentOS/RHEL
sudo yum install -y libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel
sudo yum install -y libXxf86vm-devel
```

```shell
# 1.正常打包
go build -ldflags="-s -w -H windowsgui" -o proxy-dev.exe main.go 

# 2.程序加入图标
go install fyne.io/tools/cmd/fyne@latest 
fyne package --release --id proxy.dev -os windows -icon assets/logo.png
fyne package --release --id proxy.dev -os darwin -icon assets/logo.png
fyne package --release --id proxy.dev -os linux -icon assets/logo.png

# 使用upx压缩(https://github.com/upx/upx/releases/tag/v4.1.0)
upx proxy-dev.exe
```

# 使用:

1. 安装证书
2. 开启https
3. 启动程序,查看系统代理是否生效
