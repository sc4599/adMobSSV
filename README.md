admob
###配置文件
conf中config.toml.back  
改名为：config.toml  
添加自己的广告单元ID 

### 打包
在windows下打包linux执行文件需要更改go的环境，linux系统下build则无需
```
SET  CGO_ENABLED=0
SET  GOOS=linux
SET  GOARCH=amd64
```

### 部署运行
使用编译出来的`main`到linux系统中
`chmod 777 main`
`./main`
即可运行服务，默认端口9080
可以使用supervisor 进行管理