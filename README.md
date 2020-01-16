# go-model 代码生成器

一整套的代码生成方案。

## 项目截图


## 赞赏

请作者喝杯咖啡吧！(微信号：)

![微信扫一扫]

## 特别鸣谢

- [实验楼]

## 前言

本产品致力于自动生成代码，目前生成的代码包含前后端的代码，前端采用vue框架，后端代码则使用以SpringBoot为基础的框架。

## 功能

-   [x] 服务端代码的生成。
-   [x] 前端代码的生成。
-   [x] 多种代码风格的主题。
-   [x] 支持代码的自动模式与非自动模式。
-   [x] 支持后端代码更多类型的方法。

## 使用方法

```
git clone https://github.com/yanfeng1612/go-model.git
cd go-model         // 进入目录
go build            // 项目使用golang,依赖go语言环境

./go-model d:\\go-model.json
```

## 关于配置
```
{
    "projectName" : "demo",
    "projectAdminName" : "demo-admin",
    "backTheme" : "SpringCloud",
    "frontTheme" : "Vue",
    "auto" : false,
    "firstInit" : true,
    "exportPath": "D:\\",
    "packageName" : "com.demo",
    "dataBaseName": "demo",
    "databaseIp"  : "127.0.0.1",
    "databasePort":  3306,
    "databaseUsername": "root",
    "databasePassword": "111111",
    "tableConfigList": [
        {
            "tableName": "user",
            "singleQueryList": [
                "username"
            ],
            "inListQuery": [
                "password"
            ]
        }
    ]
}
```
配置文件格式要求必须使用json,其样板如上,对配置的自动解释如下


### go-model.json
| 字段 | 用途 | 示例 | 备注 |
| ---- | ---- | --- | ----- |
| projectName | 项目英文名称 | demo | |
| projectAdminName | 前端项目名称 | demo-admin | |
| backTheme | 后端代码主题 | SpringCloud | 不同主题会生成不同风格的代码,目前支持的后端主题有SpringCloud| 
| frontTheme | 前端代码主题 | Vue        | 不同主题会生成不同风格的代码,目前支持的前端主题有Vue|
| auto | 是否主动代码生成器 | true | 支持两种形式的代码生成器,被动代码生成器-一次性生成，后续不再使用(类似与wizard),主动代码生成器(迭代使用生成器)|
| firstInit | 是否首次生成 | true | 首次生成会多生成部分脚手架类型的工具代码|
| exportPath | 生成的代码导出路径 | d:\\ |  |
| packageName | java基础包的名称 | com.demo | |
| dataBaseName | 数据库名称 | demo | |
| DatabaseIp | 数据库ip | 127.0.0.1 | |
| DatabasePort | 数据库端口 | 3306 | |
| DatabaseUsername | 数据库用户名 | root | |
| DatabasePassword | 数据库密码 | 111111| |
| tableConfigList | 自定义表中的配置 | | |
| tableConfigList.tableName | 表名称 | user | |
| tableConfigList.singleQueryList | 对表中某个字段的查询方法 | username | 对某些字段单独生成的方法,例如配置了username字段,则会生成一个方法 searchByUsername |
| tableConfigList.inListQuery | 对通用查询中增加某个字段的IN方式的查询 | username | 则在通用查询方法中会增加 usernameList字段最终传入mapper 以支持在查询中 WHERE username IN (#{usernameList}) |


## 其他注意事项

### 一、我不想编译，想直接使用，去哪下载？

## License

[MIT]
