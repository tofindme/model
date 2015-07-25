### 通过表结构来生成对应的的结构体


**生成后示例**

```
package Model

type test struct {
	id   int     `db:"id"    json:"id"`
	name varchar `db:"name"    json:"name"`
}

```


### 配置文件说明
```
{
    "dbhost": "127.0.0.1",
    "dbport": "3306",
    "user": "root",
    "pass": "root",
    "database": "information_schema",
    "schema": "test",
    "table": [
        "test",
        "employee_logins"
    ],
    "outdir": "."
}
```

**字段含义**

名称|说明|
--- | --- |
dbhost| 数据库host|
dbport| 数据库地址|
user| 数据库用户|
pass| 数据库密码|
database| 连接的数据库|
schema| 需要生成struct的表所在的schema|
table| 需要生成struct的表|
outdir| 生成struct的存放目录|


### 使用说明
- 需要安装go环境
- 在config.json文件添加相应配置参数
- 代码是在ubuntu14.04编译运行
- 使用的数据库orm是gorp使用前需下载 `go get gopkg.in/gorp.v1`
- 生成的文件代码已经通过`gofmt -s -w filename`格式化
- 准备好后可以直接`go run *.go`运行查看输出目录(outdir)下的文件




