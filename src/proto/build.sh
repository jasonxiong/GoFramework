#!/bin/sh

rm gameproto/*

cd conf/

../protoc --plugin=../protoc-gen-go --go_out=../gameproto/ *proto

