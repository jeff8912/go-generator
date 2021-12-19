#!/bin/bash

programVersion=1.0.0
compilerVersion=`go version`
buildTime=`date`
author=`whoami`
go build -ldflags "-X 'main.ProgramVersion=$programVersion' -X 'main.CompileVersion=$compilerVersion' -X 'main.BuildTime=$buildTime' -X 'main.Author=$author'"

