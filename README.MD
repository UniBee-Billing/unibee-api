# GoFrame Template For oversea-pay

Quick Start: 
- https://goframe.org/pages/viewpage.action?pageId=1114399
- goland api 自动生成配置 https://goframe.org/pages/viewpage.action?pageId=93880327
- gf gen dao https://goframe.org/pages/viewpage.action?pageId=3673173  因 数据库 default 配置问题 执行之前先删除根目录下 config.yaml 

mysql 集成
- _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
    
redis集成
- _ "github.com/gogf/gf/contrib/nosql/redis/v2"

nacos 集成
命令行参数
--nacos-ip
--nacos-port
--nacos-namespace
--nacos-group
--nacos-data-id
环境变量
nacos.ip
nacos.port
nacos.namespace
nacos.group
nacos.data.id
日常 nacos配置环境
http://daily-local.wangqubi.cn/nacos/

设置私有库： go env -w GOPRIVATE=gitlab.hknet-inc.com

打包：go build -o main main.go
临时环境运行：./main --nacos-port=7113 --nacos-namespace=hklocal --nacos-group=config --nacos-data-id=hk-go-settings.yaml --nacos-ip=frps.hknet-inc.com

[//]: # (打包上传：gf docker -tn heiku_gooverseapay:daily -p)
打包上传：sh buildAndPush.sh