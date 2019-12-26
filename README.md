# autocreate

#### 介绍
Golang 全自动生成业务框架，原来需要1天的事情，现在只需要1分钟<br/>
使用框架: gf<br/>
使用后台框架:https://github.com/CrazyRocks/goadmin<br/>
![输入图片说明](https://images.gitee.com/uploads/images/2019/1225/172404_d65f9acb_1927330.jpeg "5.jpg")<br/>
1: 生成了model<br/>
2: 生成了controller<br/>
3: 生成了router<br/>
4: 生成了moudle<br/>
5: 生成了html模式下的模板（html,js)<br/>
6: 生成了vue模式下的vue<br/>
7: 自动生成了权限menu.sql<br/>

使用方式：
1: 第一步修改数据库连接<br/>
![输入图片说明](https://images.gitee.com/uploads/images/2019/1225/172423_aaba0c76_1927330.jpeg "1.jpg")<br/>
2: 第二步创建的表和数据必须有注释(自动生成表单)<br/>
![输入图片说明](https://images.gitee.com/uploads/images/2019/1225/172435_f93a843d_1927330.jpeg "2.jpg")<br/>
3: 第三步运行 go run main.go 打开:localhost:8081<br/>
![输入图片说明](https://images.gitee.com/uploads/images/2019/1225/172449_82d6c961_1927330.jpeg "3.jpg")<br/>
4: 第四步选择表，填写项目名称，填写模块名,填写菜单名（权限菜单)<br/>
![输入图片说明](https://images.gitee.com/uploads/images/2019/1225/172458_e2a12dea_1927330.jpeg "4.jpg")<br/>
5: 第五步执行后看项目的result文件夹，直接拷贝到项目组使用<br/>

#### 软件架构
数据库: mysql<br/>
golang: latest<br/>
goframe: latest<br/>

后台管理系统两种方案可选<br/>
1:bootstrap（已经全自动生成了JS和html)<br/>
2:elementui Vue(已经全自动生成了Vue)<br/>

感谢大家支持，有钱的撒几毛，没钱的点个赞，有问题的提issue<br/>
![输入图片说明](https://images.gitee.com/uploads/images/2019/1225/172511_13869dda_1927330.jpeg "donate.jpg")

感谢大佬<br/>
| Name | Channel | Amount | Comment <br/>
|---|---|--- | ---<br/>
|李超|wechat|￥66.00 |<br/>
