#!/bin/bash
cat tag | ./crawler -u "https://book.douban.com/tag/{}?start={}&type=T" \
-s "#subject_list > ul > li > div.info > h2 > a[href] | cut / 1" \
-p @,0 -pp @,20 -np @,+20 -ep @,980 -X \
-sr 1-5,1-3 -lsr 20-50,15-30 \
-px socks5://localhost:12311 \
-ss -sc "没有找到符合条件的图书"
