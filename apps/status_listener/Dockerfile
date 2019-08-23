ARG GO_VERSION=1.11
FROM golang:${GO_VERSION}
COPY . /status_listener/
WORKDIR /status_listener/apps/status_listener
RUN go install
WORKDIR /status_listener
CMD status_listener -conf $CONF
