#!/bin/sh


gen_files=`ls gameproto/`

if [ ! -z "$gen_files" ]
then
	rm gameproto/*
fi

cd conf/

protoc --plugin=../protoc-gen-go --go_out=../gameproto/ *proto

