# autocreate

#### 介绍
Golang 全自动生成业务框架，原来需要1天的事情，现在只需要1分钟
使用框架: gf
使用后台框架:

1: 生成了model
2: 生成了controller
3: 生成了router
4: 生成了moudle
5: 生成了html模式下的模板（html,js)
6: 生成了vue模式下的vue
7: 自动生成了权限menu.sql

使用方式：
1: 第一步修改数据库连接
2: 第二步创建的表和数据必须有注释(自动生成表单)
3: 运行 go run main.go 打开:localhost:8081
4: 选择表，填写项目名称，填写模块名,填写菜单名（权限菜单)
5: 执行后看项目的result文件夹，直接拷贝到项目组使用


#### 软件架构
数据库: mysql
golang: latest
goframe: latest

后台管理系统两种方案可选
1:bootstrap（已经全自动生成了JS和html)
2:elementui Vue(已经全自动生成了Vue)

感谢大家支持，有钱的撒几毛，没钱的点个赞，有问题的提issue