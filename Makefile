gen:
	go run main.go gen

new:
	echo title of article---$$(date "+%F %T") > ./data/articles/$(F).md

new_ja:
	echo title of article---$$(date "+%F %T") > ./data/articles/ja/$(F).md
