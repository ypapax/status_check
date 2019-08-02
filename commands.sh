#!/usr/bin/env bash
set -ex

projectDir=$GOPATH/src/github.com/ypapax/status_check

build(){
	cd $projectDir/apps/status_check
	go install
}

run(){
	build
	cd $projectDir
	status_check -conf local.conf.yaml
}

runc(){
	docker-compose build
	docker-compose up
}

testl(){
	cd test
	go test -v
}
$@