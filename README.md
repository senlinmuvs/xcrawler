# xcrawler
Crawler command line tool that can crawl any website.

```
xcrawler version: 1.0.1
Options:
  -X	debug
  -cs
    	是否开启紧凑输出，多个空格和换行替换为一个
  -d string
    	保存目录,默认当前
  -ep string
    	最后一个参数
  -f string
    	从文件中取url,每行一个
  -hd string
    	设置请求头k1=v1,k2=v2
  -help
    	print help info
  -hf string
    	请求头在此文件
  -hs
    	是否横列输出
  -is string
    	pipe input separeator (default "\t")
  -js string
    	json选择器,多个逗号分隔
  -lsr string
    	100,10 表示每100秒后休息10秒
    	50-100,10-20 表示每50-100秒后息随机休息10-20秒
  -np string
    	下一页的url填充参数的变化规律,多个逗号分隔。
    	+n:表示在上一页参数基础上加n
  -p string
    	初始参数,即第一页url模版的参数,多个逗号分隔
  -pb
    	print body
  -ph
    	print response header
  -pp string
    	此参数只在有管道输入时有效，表示只有第一个管道数据使用这个初始参数，之后的都使用-p的初始参数
  -px string
    	proxy
  -q	是否打引号
  -qh
    	print request header
  -qs
    	选择器的选择项有多个时，每个结果是否打引号
  -r string
    	正则表达式,多个逗号分隔
  -rt int
    	retry count (default 3)
  -s string
    	html选择器,多个逗号分隔
  -sc string
    	停止内容：当内容包含此字符串时停止翻页
  -se
    	当json选择器内容为空时停止翻页
  -sp string
    	一个选择器对应多个结果时的分隔符 (default "\n")
  -sr string
    	10,1 表示每10次请求休息1秒
    	1-10,1-5 表示每1-10次请休息随机1-5秒
  -ss
    	当此页内容与上一页内容相同时停止翻页
  -to int
    	timeout seconds (default 10)
  -tp int
    	总共翻多少页
  -u string
    	url或起始url模版
  -v	show version
```