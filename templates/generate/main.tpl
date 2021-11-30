package main

import (
    "{{.packagePrefix}}/common"
    "{{.packagePrefix}}/config"
    "{{.packagePrefix}}/control"
    "{{.packagePrefix}}/log"
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
