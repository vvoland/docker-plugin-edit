.PHONY: binary
.PHONY: vendor

binary:
	docker buildx build . -o type=local,dest=bin

install: binary
	cp bin/docker-edit ~/.docker/cli-plugins/docker-edit

vendor:
	docker run --init -it --rm \
		-v docker-plugin-edit-cache:/root/.cache \
		-v ./:/proj -w /proj \
		$(shell docker buildx build . -q --target base) \
			go mod tidy
