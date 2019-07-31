#!/usr/bin/env bash
set -ex

build(){
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