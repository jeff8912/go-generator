package main

import (
	"go-generator/common"
	"go-generator/config"
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
	common.BannerShow(ProgramVersion, CompileVersion, BuildTime, Author)

	log.Init(common.GetMainName(), config.GetValue("log", "level"))
	defer log.Sync()

	control.MainControl()
}
