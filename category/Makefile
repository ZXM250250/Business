GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	sudo docker run --rm  -v $(shell pwd):$(shell pwd) -w $(shell pwd) -e ICODE=CF388DF1EF1C5EBE cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/category/category.proto

.PHONY: build
build: 

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o category-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t category-service:latest
