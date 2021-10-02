serve: 
	go run *.go
gen:
	go run *.go -gen
new:
	go run *.go -new

deploy: gen
	git commit -m 'deploy' -a
	git push origin master
