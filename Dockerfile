FROM golang:1.17 as devimage
ENV GOLANG_CI_LINT_VERSION=v1.18.0
ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.neoway.com.br
ENV GOSUMDB=off
RUN cd /usr && \
    wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s ${GOLANG_CI_LINT_VERSION}
RUN go get golang.org/x/perf/cmd/benchstat
WORKDIR /app
COPY go.mod /app
RUN go mod download
RUN go mod tidy
COPY . /app
EXPOSE 80

FROM devimage as buildimage
ARG version
ENV CGO_ENABLED=0
RUN go build -a -installsuffix cgo -ldflags "-w -s -X main.Version=$version" -o ./cmd/semantic-release/semantic-release ./cmd/semantic-release

FROM alpine:3.9 as prodimage
COPY --from=buildimage /app/cmd/semantic-release/semantic-release /app/semantic-release
EXPOSE 80
ENTRYPOINT ["/app/semantic-release"]