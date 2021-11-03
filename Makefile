serve: 
	go run *.go
gen:
	go run *.go -gen

deploy: gen
	git add .
	git commit -m 'deploy'
	git push origin master
