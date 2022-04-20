FROM alpine:3.9 as base

RUN apk --no-cache update && \
    apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

RUN adduser -D -g '' appuser

COPY ./cmd/semantic-release/semantic-release /app/semantic-release
ARG SSH_INTEGRATION_SEMANTIC
ARG GITLAB_CONTAINER_NAME
ENV SSH_INTEGRATION_SEMANTIC=$SSH_INTEGRATION_SEMANTIC
ENV GITLAB_CONTAINER_NAME=$GITLAB_CONTAINER_NAME
RUN echo $SSH_INTEGRATION_SEMANTIC
RUN echo "172.20.0.2 gitlab.integration-tests.com" >> /etc/hosts


FROM scratch

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /app/semantic-release /app/semantic-release

# Use an unprivileged user.
USER appuser

ENTRYPOINT ["/app/semantic-release"]
