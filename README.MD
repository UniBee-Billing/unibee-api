# GoFrame Template For oversea-pay

#### 日常开发环境地址：https://daily.hknet-inc.com/gooverseapay/swagger
##### Mock-Controller 模拟 OpenApi 测试使用
##### Open-Controller OpenApi 接口-需 token 调用

### 本地开发
#### 设置私有库： go env -w GOPRIVATE=gitlab.hknet-inc.com

#### 本地环境运行，添加Program Arguments: --nacos-port=30099 --nacos-namespace=local --nacos-group=config --nacos-data-id=unib-settings.yaml --nacos-ip=43.207.114.88
##### go run main.go --nacos-port=30099 --nacos-namespace=local --nacos-group=config --nacos-data-id=unib-settings.yaml --nacos-ip=43.207.114.88

[//]: # (打包上传：gf docker -tn heiku_gooverseapay:daily -p)
#### 生成 api 代码：sh generate_ctrl.sh
#### 生成 dao 代码：sh generale_dao.sh
#### 打包上传：sh buildAndPush.sh

#### AWS 开发环境
https://unibee.atlassian.net/wiki/spaces/SD/pages/4292713/AWS+Develop+Environment


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