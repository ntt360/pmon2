#!/usr/bin/env bash

#--------------------------------
# pmon2开发初始化初始化开发环境
#--------------------------------

set -e

rootDir=$(cd "$(dirname "$0")"; pwd)

sudo rm -rf "$rootDir/config/config-dev.yml" "$rootDir/tmp"

# 写配置
cd "config"
echo "data: $rootDir/tmp/data" >> "config-dev.yml"
echo "logs: $rootDir/tmp/logs" >> "config-dev.yml"
cd ..

#创建文件夹
logs="$rootDir/tmp/logs"
data="$rootDir/tmp/data"
if [ -d "$data" ]; then
  rm -rf "$data"
fi
mkdir -p "$data"

if [ -d "$logs" ]; then
  rm -rf "$logs"
fi
mkdir -p "$logs"

# build go bin
go mod tidy
go build -o bin/pmon2 cmd/pmon2/pmon2.go
go build -o bin/pmond cmd/pmond/pmond.go
# `bin/test` 进程是用于测试模拟的业务进程，开发也可以测试自己的业务进程。
go build -o bin/test test/test.go

