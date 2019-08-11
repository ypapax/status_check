ARG GO_VERSION=1.11
FROM golang:${GO_VERSION}
COPY . /api/
WORKDIR /api/apps/api
RUN go install
WORKDIR /api
CMD api -conf $CONF
