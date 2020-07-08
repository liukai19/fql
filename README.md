# FQL
基于 SQL的数据文件分析统计工具
### 背景
- 基于 AWK/SORT/UNIQ/GREP 等 Shell 命令实现日志分析和统计(Ng/Apache/PHP),复杂的 AWK 隐式语法,多种命令入参,长时间不用很难熟记。
- 产品&运营等基于 Execl 分析数据，本身 Execl 功能强大，但是没有系统的学习,比如数据透视等功能，很难对数据分析

### 介绍
FQL是一个简单的行列数据文件分析工具，基于 SQL(结构化查询语言)实现各种分析统计类场景，覆盖 AWK/SORT/UNIQ/GREP等绝大部分分析统计功能。

#### 简要理论介绍
FQL有以下4部分构成:
- 记号提取器
- Query切分器
- 抽象语法树
- 算子管理器

#### 安装
```
git clone https://github.com/xiaoxuz/fql.git
echo 'FQL_HOME=xxx' >> ~/.bashrc
echo 'export PATH=$PATH:$FQL_HOME/' >> ~/.bashrc
source ~/.bashrc
```
#### 使用
```
FQL Version: fql/1.0.0
Usage: fql [-h] [-s sql] [-d delim]

Options:
  -d string
    	set the column delimiter, the default Spaces (default " ")
  -h	this help
  -s SQL
    	send an SQL parse data file
```

#### SQL 基础结构
```
SQL = SELECT 
[field, func ...] 
[FROM file] 
[WHERE {where_condition ...} {AND | OR}] 
[GROUP BY {col_name, ...} 
[ORDER BY {col_name} {ASC | DESC}] 
[LIMIT {[offset,] count}]
```
#### 列名标识
`@{1,2,3,4...N}`

e.g : @1 第一列
### 条件运算符
- `>`
- `<`
- `=`
- `!=`
- `>=`
- `<=`
- `IN ([v, ...])`
- `NOT IN ([v, ...])`
- `LIKE 'xxx'`
#### 聚合函数
- SUM : 返回数值列的总数（`float64`) 
    - e.g : `SELECT SUM(column_name) FROM file_name`
- MAX : 返回一列中的最大值 (`float64`) 
    - e.g : `SELECT MAX(column_name) FROM file_name`
- MIN : 返回一列中的最小值 (`float64`) 
    - e.g : `SELECT MIN(column_name) FROM file_name`
- AVG : 返回数值列的平均值 (`float64`) 
    - e.g : `SELECT AVG(column_name) FROM file_name`
- COUNT : 返回匹配指定条件的行数 (`int64`) 
    - e.g : `SELECT COUNT(column_name) FROM file_name`
- DISTINCT : 返回唯一不同的值 
    - e.g : `SELECT DISTINCT(column_name) FROM file_name`
- FROM_UNIXTIME :  将时间戳转化成指定模块的时间 
    - e.g : `SELECT FROM_UNIXTIME(column_name::'%Y-%m-%d %H:%i:%s') FROM file_name`
- UNIX_TIMESTAMP : 将时间字符串按照指定格式解析成时间戳
    - e.g : `SELECT UNIX_TIMESTAMP(column_name::'[%d/%m/%y:%h:%i:%s') FROM file_name`

### For Example
Ng access.log demo
```
21.10.134.3 - - [03/Jul/2020:17:41:29 +0800] "GET /ttt/l/getinfo?lId=6421786 HTTP/1.0" 200 2290 "https://dddd.test.cc/ssss/view/ttt/task/call-out/contact" "S=IPS_2284aebaa3f5d964a8e0930059af1898" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36" 0.059 2489926140 21.10.134.3 21.10.128.156 unix:/hhhh/hhhhwork/var/php-cgi.sock dddd.test.cc "12.17.39.190" hhhhwork question 24899261402625677504070317 1593769289.982 0.786
```
---
查看访问时间戳
```
fql -s "select @1 as ip, UNIX_TIMESTAMP(@4::'[%d/%em/%y:%h:%i:%s') as ts, @4 ,@7 as uri from test/access.log limit 10"

IP | TS | @4 | URI
21.10.128.229 | 1593766881 | [03/Jul/2020:17:01:21 | /llllcall/api/getrecord
21.10.128.229 | 1593766881 | [03/Jul/2020:17:01:21 | /llllcall/api/getrecord
21.10.128.161 | 1593766881 | [03/Jul/2020:17:01:21 | /llllcall/api/getrecord
21.10.134.3 | 1593766881 | [03/Jul/2020:17:01:21 | /aaasc/leads/editname?courseId=567809&leadsId=106951695&name=%E5%BC%A0%E9%9B%A8%E6%B3%BD
21.10.134.4 | 1593766881 | [03/Jul/2020:17:01:21 | /aaasc/call/callout?courseId=564769&leadsId=106309712
21.10.130.197 | 1593766881 | [03/Jul/2020:17:01:21 | /datacenter/datacenter/api/search
21.10.130.197 | 1593766881 | [03/Jul/2020:17:01:21 | /callcenter/auto/existagentlist
21.10.48.16 | 1593766881 | [03/Jul/2020:17:01:21 | /llllcore/commit/commit?channel=dal&cmdno=810003&topic=core&transid=2542700024
21.10.130.204 | 1593766881 | [03/Jul/2020:17:01:21 | /callcenter/auto/existagentlist
21.10.134.77 | 1593766881 | [03/Jul/2020:17:01:21 | /datacenter/datacenter/api/search
```
---

统计访问 PV/UV
```
fql -s "select count(@1) as PV, count(distinct(@1)) as UV from test/access.log"

PV | UV
77394 | 150
```
---
统计访问频率最高的 Top10 Uri
```
fql -s "select @7 as uri, count(@7) as num from test/access.log group by @7 order by num desc limit 10"

URI | NUM
/llllcall/api/getrecord | 6743
/data/data/api/search | 4375
/ttt/le/getle| 1684
/aaasc/api/getst| 1657
/llllcall/api/getcollectlist | 1587
/perfo/api/get| 1577
/ins/api/check| 1332
/ins/api/getmaterial| 1325
/ins/api/getdist| 1254
/monitor.php | 1212
```
---
 统计某一接口平均耗时
```
fql -s "select avg(@27) as avg from test/access.log where @7 = '/llllcall/api/getrecord'"

AVG
1.30
```

### TODO
- UnitTest补充
- 嵌套查询、子查询
- 多文件 JOIN 查询

### About
![avatar](https://github.com/xiaoxuz/fql/blob/master/wechat.jpg)