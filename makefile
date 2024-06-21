.PHONY: build
.PHONY: run

build:
	go build -o ./kubestream

run:
	go run main.go