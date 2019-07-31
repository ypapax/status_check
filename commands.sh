#!/usr/bin/env bash
set -ex

build(){
	cd $GOPATH/src/github.com/ypapax/status_check/apps/status_check
	go install
}

run(){
	build
	status_check --database psql
}

runc(){
	docker-compose build
	docker-compose up
}

$@