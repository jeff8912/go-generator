package common

import (
	"flag"
	"fmt"
	"github.com/dimiro1/banner"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetMainName() string {
	return strings.Split(filepath.Base(os.Args[0]), ".")[0]
}

func GetMainPath() string {
	fileAbsPath, _ := filepath.Abs(os.Args[0])
	programPath := filepath.Dir(fileAbsPath)

	return programPath
}

func GetWorkPath() string {
	workPath, _ := os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	return workPath
}

func JoinPath(arrStr ...string) string {
	return strings.Join(arrStr, "/")
}

func GetTcpAddress(ip string, port interface{}) string {
	var strPort string
	switch value := port.(type) {
	case int:
		strPort = strconv.Itoa(value)
	case string:
		strPort = value
	}
	return ip + ":" + strPort
}

// Convert uint to string
func InetNtoa(ipnr int64) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

// Convert string to int64
func InetAton(ipnr string) int64 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) {{$.lt}}{{$.lt}} 24
	sum += int64(b1) {{$.lt}}{{$.lt}} 16
	sum += int64(b2) {{$.lt}}{{$.lt}} 8
	sum += int64(b3)

	return sum
}



func BannerShow(programVersion, compilerVersion, buildTime, author string)  {
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
***           ."" '\ '.___\_(|)_/___.'  /'"".             ***
***         | | :  '- \'.;'\ _ /';.'/ - ' : | |            ***
***         \  \ '-.   \_ __\ /__ _/   .-' /  /            ***
***   ========'-.____'-.___\_____/___.-'____.-'========    ***
***                        '=---='                         ***
***   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^   ***
***         佛祖保佑       永不宕机     永无BUG                ***
**************************************************************
****************** Compile Environment ***********************
Program version : %s
Compiler version : %s
Build time : %s
Author : %s
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