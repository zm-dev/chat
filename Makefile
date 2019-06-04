docker:
	docker build -t chat:latest -f Dockerfile.server .
docker_upload: docker
	docker tag chat:latest registry.cn-hangzhou.aliyuncs.com/zm-dev/chat:latest
	docker push registry.cn-hangzhou.aliyuncs.com/zm-dev/chat:latest
run:
	docker run -d chat:latest
