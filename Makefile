GOPATH?=$(HOME)/go

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/sadwave-events-api-v2 ./main.go

docker-build:
	docker build -t punxlab/sadwave-events-api-v2 .

docker-run:
	docker run punxlab/sadwave-events-api-v2:latest

docker-push:
	docker push punxlab/sadwave-events-api-v2:latest