## 后端部署方法


### 1. 启动mysql
- 按照config.json里的端口、用户名和密码，配置本地mysql数据库（注意用户要有建表权限）
- 在后端机器上启动mysql服务
- 在后端服务启动并成功链接后，将自动建表

### 2. 前后端配置Https公钥
- 向（阿里云等服务商）申请SSL证书，成功后生成密钥，其中key和pem文件分别命名为hita.key/hita.pem放置在后端代码的middelware/下，而bks文件则同样命名为hita.bks放置在Android代码的component/src/main/res/raw
- Android端需要在SslContextFactory中配置信任证书的密码，然后取消BaseWebSource.kt第15～16行的注释

### 3. 启动后端程序
- 安装go
```shell
sudo apt-get update && sudo apt-get -y install golang-go 
```
- 进入项目，启动后端
```shell
cd HITA_Go # 进入到项目目录
go build main.go # 编译程序
./main # 启动程序
```