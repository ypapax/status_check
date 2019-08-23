#!/usr/bin/env bash
set -ex

projectDir=$GOPATH/src/github.com/ypapax/status_check

build(){
  pushd $projectDir/apps/api
	go install
	popd

	pushd $projectDir/apps/status_check
	go install
	popd

	pushd $projectDir/apps/status_listener
	go install
	popd
}

run(){
	build
	cd $projectDir
	status_check -conf local.conf.yaml
}

runc(){
  build
	set +e; docker network create status-check-network; set -e;
	docker-compose build
	docker-compose down
	docker-compose up
}

avail(){
  to=$(date +%s)
  from=$((to-600))
  curl localhost:3000/services-count/available/$from/$to
}

test(){
  build
  cd queue
  go test -v
	cd ../test
	go test -v -timeout 20m
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