package main

import (
    "{{.packagePrefix}}/banner"
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
	defer log.Sync()
	banner.Show(ProgramVersion, CompileVersion, BuildTime, Author)
	control.MainControl()
}
