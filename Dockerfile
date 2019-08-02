ARG GO_VERSION=1.11
FROM golang:${GO_VERSION}
COPY . /status_check/

WORKDIR /status_check/apps/status_check
RUN go install
WORKDIR /status_check
RUN chmod +x /status_check/entrypoint.sh
ENTRYPOINT "/status_check/entrypoint.sh"
