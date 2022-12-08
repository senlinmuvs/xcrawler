#!/bin/bash
cat m_id | ./crawler -u https://movie.douban.com/subject/{} \
-s "#content > h1 > span:nth-child(1),#info,#link-report-intra > span,#interest_sectl strong.rating_num,#mainpic > a > img[src],#info > span:nth-child(1) > span.attrs > a[href]|cut / 1,#info > span:nth-child(3) > span.attrs > a[href]|cut / 1" \
-hs -sp ";" -q -se -cs -X \
-sr 2-5,1-3 -lsr 20-40,15-30 \
-px socks5://localhost:12311