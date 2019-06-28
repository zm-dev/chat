upload:
	docker build -t chat:latest -f Dockerfile.dev .
	docker tag chat:latest registry.cn-hangzhou.aliyuncs.com/zm-dev/chat:latest
	docker push registry.cn-hangzhou.aliyuncs.com/zm-dev/chat:latest

release:
	docker build -t chat:release -f Dockerfile.release .
	docker tag chat:release registry.cn-hangzhou.aliyuncs.com/zm-dev/chat:release
	docker push registry.cn-hangzhou.aliyuncs.com/zm-dev/chat:release

run:
	docker run -d chat:latest

all: upload release