serve: 
	go run *.go
gen:
	go run *.go -gen
new:
	go run *.go -new

deploy: gen
	git add .
	git commit -m 'deploy'
	git push origin master
