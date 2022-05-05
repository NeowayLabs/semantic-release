version ?= latest
img = neowaylabs/semantic-release:$(version)
imgdev = neowaylabs/semantic-releasedev:$(version)
gitlab_container_name=gitlab_integration_tests
uid=$(shell id -u $$USER)
gid=$(shell id -g $$USER)
dockerbuilduser=--build-arg USER_ID=$(uid) --build-arg GROUP_ID=$(gid) --build-arg USER
wd=$(shell pwd)
modcachedir=$(wd)/.gomodcachedir
cachevol=$(modcachedir):/go/pkg/mod
appvol=$(wd):/app
run=docker run --rm -v $(appvol) -v $(cachevol) $(imgdev)
runbuild=docker run --rm -ti -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 -v $(appvol) -v $(cachevol) $(imgdev)
cov=coverage.out
covhtml=coverage.html
git_group_test=dataplatform
git_project_test=integration-tests
git_host_test=gitlab.integration-tests.com
run_local=./cmd/semantic-release/semantic-release

all: check build

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

# WHY: If cache dir does not exist it is mapped inside container as root
# If it exists it is mapped belonging to the non-root user inside the container
modcache:
	@mkdir -p $(modcachedir)

image: build
	docker build . -t $(img)

imagedev:
	docker build . -t $(imgdev) -f ./hack/Dockerfile $(dockerbuilduser) --build-arg SSH_INTEGRATION_SEMANTIC="${SSH_INTEGRATION_SEMANTIC}"

remove-images:
	docker rmi -f $(img) && docker rmi -f $(imgdev)

release: guard-version publish
	git tag -a $(version) -m "Generated release "$(version)
	git push origin $(version)

publish: image
	docker push $(img)

build: modcache imagedev
	$(runbuild) go build -v -ldflags "-w -s -X main.Version=$(version)" -o $(run_local) ./cmd/semantic-release

env: ##@environment Create network and run gitlab container.
	GITLAB_CONTAINER_NAME=${gitlab_container_name} \
	docker-compose up --build -d

shell-git:
	docker exec -it ${gitlab_container_name} /bin/bash

log-git:
	docker logs -f ${gitlab_container_name}

git-ip:
	docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${gitlab_container_name}

inspect-git:
	docker inspect ${gitlab_container_name}

create-gitlab-backup:
	./hack/gitlab-backup.sh create ${gitlab_container_name}

restore-gitlab-backup:
	./hack/gitlab-backup.sh restore ${gitlab_container_name}

run-gitlab-env: env restore-gitlab-backup

env-stop: ##@environment Remove gitlab container, and network.
	GITLAB_CONTAINER_NAME=${gitlab_container_name} docker-compose kill
	GITLAB_CONTAINER_NAME=${gitlab_container_name} docker-compose rm -v -f

clean-go-build:
	./hack/clean-go-build.sh

build-go: clean-go-build
	cd cmd/semantic-release && go build -o semantic-release

run-local: build-go
	$(run_local) up -git-group ${git_group_test} -git-project ${git_project_test} -git-host ${git_host_test} -auth "${SSH_INTEGRATION_SEMANTIC}" -setup-py true

run-help:
	$(run_local) help

run-help-cmt:
	$(run_local) help-cmt

run-docker-local:
	docker run ${img} up -git-group ${git_group_test} -git-project ${git_project_test} -git-host ${git_host_test} -setup-py true -auth ${SSH_INTEGRATION_SEMANTIC}

check: modcache imagedev
	$(run) go test -tags unit -timeout 20s -race -coverprofile=$(cov) ./...

check-integration: run-gitlab-env image
	./hack/check-integration.sh $(parameters)

coverage: modcache check
	$(run) go tool cover -html=$(cov) -o=$(covhtml)
	xdg-open coverage.html

static-analysis: modcache imagedev
	$(run) golangci-lint run ./...

modtidy: modcache imagedev
	$(run) go mod tidy

fmt: modcache imagedev
	$(run) gofmt -w -s -l .

githooks:
	@echo "copying git hooks"
	@mkdir -p .git/hooks
	@cp hack/githooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "git hooks copied"

shell: modcache imagedev
	$(run) sh