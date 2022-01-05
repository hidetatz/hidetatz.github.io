gen: 
	go run *.go
serve:
	go run *.go -serve

deploy: gen
	git add .
	git commit -m 'deploy'
	git push origin master
