DOCKER_IMG="filter-bot:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-img:
	docker build \
		--no-cache \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f $(CURDIR)/Dockerfile .

run-img: build-img
	docker run --env-file $(CURDIR)/.env $(DOCKER_IMG)

.PHONY: build-img run-img