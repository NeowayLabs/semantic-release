FROM golang:1.16-stretch

ENV GOLANG_CI_LINT_VERSION=v1.13.2

RUN cd /usr && \
    wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s ${GOLANG_CI_LINT_VERSION}

ARG USER
ARG USER_ID
ARG GROUP_ID
ARG GITLAB_CONTAINER_NAME
ENV GITLAB_CONTAINER_NAME=$GITLAB_CONTAINER_NAME
RUN echo "172.20.0.2 gitlab.integration-tests.com" >> /etc/hosts

RUN groupadd -f -g ${GROUP_ID} ${USER} && \
    useradd -m -g ${GROUP_ID} -u ${USER_ID} ${USER} || echo "user already exists"

USER ${USER_ID}:${GROUP_ID}

WORKDIR /app
