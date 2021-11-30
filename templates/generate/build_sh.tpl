#!/bin/bash

program_version="1.0.3"
compiler_version=`go version`
build_time=`date`
author=`whoami`
go build -ldflags "-X 'main.PROGRAM_VERSION=$program_version' -X 'main.COMPILER_VERSION=$compiler_version' -X 'main.BUILD_TIME=$build_time' -X 'main.AUTHOR=$author'"
