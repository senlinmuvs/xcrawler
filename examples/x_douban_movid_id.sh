#/bin/bash
cat movie_tags|./crawler -u "https://m.douban.com/rexxar/api/v2/movie/recommend?refresh=0&start={}&count=20&selected_categories=%7B%7D&uncollect=false&sort=R&tags={}" \
-hf header.txt -js items.id \
-p 0,@ -np +20,@ -se -X \
-sr 5-15,1-3 -lsr 40-80,15-30