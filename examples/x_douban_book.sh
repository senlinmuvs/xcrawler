#!/bin/bash
cat ids | ./crawler \
-u "https://book.douban.com/subject/{}/" \
-s "#wrapper > h1 > span,#info,div.related_info span.all div.intro,#content > div > div.article > div.related_info > div:nth-child(5) > div > div,#interest_sectl strong.rating_num,#mainpic > a > img[src]" \
-X -sr 1-5,1-3 -lsr 20-50,15-30 \
-ss -sc "没有找到符合条件的图书" -hs -cs -q
#-px socks5://localhost:12311 \
