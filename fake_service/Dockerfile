ARG GO_VERSION=1.11
FROM golang:${GO_VERSION}
COPY . /fake_service/

WORKDIR /fake_service/apps/fake_services
RUN go install
WORKDIR /fake_service
RUN chmod +x /fake_service/fake_service/entrypoint.sh
ENTRYPOINT "/fake_service/fake_service/entrypoint.sh"
