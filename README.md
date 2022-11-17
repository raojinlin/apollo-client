## README

这是一个Apollo配置中心的客户端，可以为应用提供语言无关的配置文件，使用.env格式。

### 特性
1. 语言无关
2. 支持dotenv
3. 配置推送

### 使用

命令参数
```shell
A fast and easy to use command line tool for apollo

Usage:
  apollo-client [flags]

Flags:
  -i, --appId string            Apollo appId
  -d, --output string           Path to save config. (default "./")
  -c, --cluster string          Cluster name (default "default")
      --config string           config file (default "./config.yaml")
  -h, --help                    help for apollo-client
      --namespace stringArray   App namespace
  -n, --notify string           notify command while config changed.
  -u, --notifyUrl string        Push to server if config changed.
  -s, --server string           Apollo config server address eg. http://192.168.31.111:8081/
      --viper                   Use Viper for configuration (default true)
      --watch                   Listen for configuration change
```


执行命令
```shell
# 不使用配置文件
./apollo-client --server http://apollo-dev.xxxx.cc:8082\
 --cluster default --appId application --cacheDir /tmp
# 使用配置文件，默认使用
./apollo-client --config ./config.yaml
```

使用配置文件，默认是./config.yaml

```yaml
# 服务器地址
server: http://apollo-dev.xxxx.cc:8082
# 是否监听服务器配置更新
watch: False
appId: node-app
# 配置保留路径
output: ./
# 所属集群
cluster: default
# 当配置变更时，notify脚本内容
notify: |
  env
  echo updated
# 命名空间
namespaces:
  - application
```

命令执行完成后会在```cacheDir```看到```application.env```文件生成。
```dotenv
AWS_REGION="cn-northwest-1"
AWS_S3_ACCOUNT_ID="123123123"
AWS_S3_ACCESS_KEY_ID="xxxxxxx"
AWS_S3_REGION="ap-southeast-1"
```

### 监听配置变更
当指定了```watch```参数后会监听服务器的配置变更。指定```notify```参数可以填写通知脚本，程序会自动生成一个脚本文件，并在更新后执行它。
```shell
$ ./apollo-client --watch --notify 'echo update'
```
传入通知脚本的环境变量:

| Environment Variable | Description                   |
|----------------------|-------------------------------|
| NAMESPACES           | 传入的namespace，多个namespace用,分隔 |
| CLUSTER              | 当前集群                          |
| CACHEDIR             | 配置文件的保存路径                     |
| APPID                | App ID                        |
| SERVER               | Apollo服务器地址                   |


### 监听配置更改并推送到服务器
指定```notifiyUrl```参数，将最新配置推送到服务器。
```shell
$ ./apollo-client --watch --notifyUrl http://127.0.0.1/env/notify
```