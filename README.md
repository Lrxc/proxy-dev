proxy-dev

GUI框架: wails

系统级抓包, 请求重定向,拦截,修改等

# 功能列表:

- 网络拦截重定向

# 打包

```shell
go install fyne.io/tools/cmd/fyne@latest # 安装 fyne cmd
fyne package --release --id proxy.dev -os windows -icon assets/logo.png # windows加入图标打包
```

包太大,剔除多余并压缩

```bash
#最小打包
#-ldflags=“参数”： 表示将引号里面的参数传给编译器
#-s：去掉符号信息（这样panic时，stack trace就没有任何文件名/行号信息了，这等价于普通C/C+=程序被strip的效果）
#-w：去掉DWARF调试信息 （得到的程序就不能用gdb调试了）
#-H windowsgui : 以windows gui形式打包，不带dos窗口。其中注意H是大写的
go build -ldflags="-s -w -H windowsgui" -o proxy-dev.exe main.go 

#使用upx再次压缩(https://github.com/upx/upx/releases/tag/v4.1.0)
upx -9 proxy-dev.exe
```

# 使用:

1. 安装证书
2. 开启https
3. 启动程序,查看系统代理是否生效
