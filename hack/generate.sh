#!/bin/bash

proto(){
	protoc --go_out=. ./pkg/api/*.proto
}


main(){
	case $1 in
		"proto" ) proto
		;;
	esac
}

main $@





