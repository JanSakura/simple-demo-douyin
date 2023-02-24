# simple-demo-douyin

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

工程lib文件已添加

```shell
go run ./main.go
```

### 功能说明

* 系统入口main.go
* routers里InitRouter加载路由和SQL数据库
* Init初始化MySQL数据库，其中配置文件采用toml格式
* 各路由启动对应的Handler， 并做中间件检测
* Handler调用对应的service功能，对结果做判断
* service调用dao ，有需要的使用Redis做缓存，对得到的数据进行操作；FFmpeg文件存在lib
* dao抽象并返回对应models的CRUD结果
* cache初始化Redis，并实现Redis的数据库操作
* 上传的文件存在static里
