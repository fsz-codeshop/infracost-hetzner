BINARY_NAME=infracost-hetzner
DOCKER_IMAGE=fsz-codeshop/infracost-hetzner

build-local:
	go build -o $(BINARY_NAME) main.go

build-docker:
	docker build -t $(DOCKER_IMAGE) .

run-docker:
	docker run --rm \
		-e HCLOUD_TOKEN=$(HCLOUD_TOKEN) \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e GITHUB_REPOSITORY=$(GITHUB_REPOSITORY) \
		-e PR_NUMBER=$(PR_NUMBER) \
		-v $(PWD):/app \
		$(DOCKER_IMAGE) --plan /app/plan.json

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
