# GoFrame Template For oversea-pay

#### 设置私有库： go env -w GOPRIVATE=gitlab.hknet-inc.com

#### 本地环境运行，添加Program Arguments: --nacos-port=7113 --nacos-namespace=hklocal --nacos-group=config --nacos-data-id=hk-go-settings.yaml --nacos-ip=frps.hknet-inc.com

[//]: # (打包上传：gf docker -tn heiku_gooverseapay:daily -p)
#### 生成 api 代码：sh generate_ctrl.sh
#### 生成 dao 代码：sh generale_dao.sh
#### 打包上传：sh buildAndPush.sh

#### 数据库账号
#### 地址：frps.hknet-inc.com
#### 端口：7111
#### 账号：gooverseapay
#### 密码：1lc3BhY2UiOiJrdWJlcm5
#### create user 'gooverseapay'@'%' identified with mysql_native_password by '1lc3BhY2UiOiJrdWJlcm5';
#### GRANT ALL ON gooversea_pay.* TO 'gooverseapay'@'%';
#### GRANT ALL ON xin.* TO 'gooverseapay'@'%';
#### FLUSH PRIVILEGES;


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