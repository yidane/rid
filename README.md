# rid
help dev sync database


```
output [dir]  设置输出目录
load [database] 指定名称，加载该数据库所有表；不指定，加载所有数据库信息
use [database] 使用某数据库
add [table] 向缓冲区添加表，* 表示下载该数据库下所有表
rm [table] 移除缓冲区内的表
list 展示缓冲区内所有表
clear 清除缓冲区
download -r=true 下载缓冲区内所有表数据,-r表示若已经下载过，是否重新下载
```