#!/usr/bin/env bash
set -ex

projectDir=$GOPATH/src/github.com/ypapax/status_check

build(){
  cd $projectDir/apps/api
	go install
	cd $projectDir/apps/status_check
	go install
}

run(){
	build
	cd $projectDir
	status_check -conf local.conf.yaml
}

runc(){
	set +e; docker network create status-check-network; set -e;
	docker-compose build
	docker-compose up
}

avail(){
  to=$(date +%s)
  from=$((to-3600))
  curl localhost:3000/services-count/available/$from/$to
}

test(){
  cd queue
  go test -v
	cd ../test
	go test -v
}

testYamlFile=./docker-compose-test.yml
#test(){
#	docker-compose -f $testYamlFile build
#	docker-compose -f $testYamlFile up
#}

bf(){
	docker-compose -f $testYamlFile build fake-service
	docker-compose -f $testYamlFile up fake-service
}

"$@"