IMAGE_NAME ?= matburt/minifluxauto:latest

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/minifluxauto

docker: build
	docker build -t $(IMAGE_NAME) .