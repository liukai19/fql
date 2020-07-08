# FQL
基于 SQL的数据文件分析统计工具
### 背景
- 基于 AWK/SORT/UNIQ/GREP 等 Shell 命令实现日志分析和统计(Ng/Apache/PHP),复杂的 AWK 隐式语法,多种命令入参,长时间不用很难熟记。
- 产品&运营等基于 Execl 分析数据，本身 Execl 功能强大，但是没有系统的学习,比如数据透视等功能，很难对数据分析

### 介绍
FQL是一个简单的行列数据文件分析工具，基于 SQL(结构化查询语言)实现各种分析统计类场景，覆盖 AWK/SORT/UNIQ/GREP等绝大部分分析统计功能。

#### 简要理论介绍
FQL有已下4部分构成:
- 记号提取器
- Query切分器
- 抽象语法树
- 算子管理器

#### 安装
```
git clone https://github.com/xiaoxuz/fql.git
echo 'export FQL_HOME="xxx"' >> ~/.bashrc
echo 'export PATH="$PATH:$FQL_HOME/' >> ~/.bashrc
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
