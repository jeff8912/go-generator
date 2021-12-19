package banner

import (
	"flag"
	"fmt"
	"github.com/dimiro1/banner"
	"os"
	"strings"
)

func Show(programVersion, compilerVersion, buildTime, author string) {
	bannerLogo :=
		`
**************************************************************
***                       _ooOoo_                          ***
***                      o8888888o                         ***
***                      88" . "88                         ***
***                      (| ^_^ |)                         ***
***                      O\  =  /O                         ***
***                   ____/'---'\____                      ***
***                 .'  \\|     |//  '.                    ***
***                /  \\|||  :  |||//  \                   ***
***               /  _||||| -:- |||||-  \                  ***
***               |   | \\\  -  /// |   |                  ***
***               | \_|  ''\---/''  |   |                  ***
***               \  .-\__  '-'  ___/-. /                  ***
***             ___'. .'  /--.--\  '. . ___                ***
***           ."" '\ '.___\_(|)_/___.'  /'"".              ***
***         | | :  '- \'.;'\ _ /';.'/ - ' : | |            ***
***         \  \ '-.   \_ __\ /__ _/   .-' /  /            ***
***   ========'-.____'-.___\_____/___.-'____.-'========    ***
***                        '=---='                         ***
***   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^   ***
***         佛祖保佑       永不宕机     永无BUG             ***
**************************************************************
****************** Compile Environment ***********************
Program version : %s
Compiler version : %s
Build time : %s
Author : %s
**************************************************************
****************** Running Environment ***********************
Go running version : {{ .GoVersion }}
Go running OS : {{ .GOOS }}
Startup time : {{ .Now "2006-01-02 15:04:05" }}
**************************************************************

`
	var version bool
	flag.BoolVar(&version, "v", false, "print the version info")
	flag.Parse()

	newBanner := fmt.Sprintf(bannerLogo, programVersion, compilerVersion, buildTime, author)

	banner.Init(os.Stdout, true, true, strings.NewReader(newBanner))

	if version {
		os.Exit(0)
	}
}
