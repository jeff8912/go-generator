package main

import (
	"go-generator/banner"
	"go-generator/control"
	"go-generator/log"
)

var (
	ProgramVersion string
	CompileVersion string
	BuildTime      string
	Author         string
)

func main() {
	defer log.Sync()
	banner.Show(ProgramVersion, CompileVersion, BuildTime, Author)
	control.MainControl()
}
