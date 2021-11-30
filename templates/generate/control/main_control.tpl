package control

import (
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"os/signal"
	"{{.packagePrefix}}/handlers"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/context"

	"{{.packagePrefix}}/common"
	"{{.packagePrefix}}/config"
	"{{.packagePrefix}}/log"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echoLog "github.com/labstack/gommon/log"
)


type (
	handlerEntity struct {
		handlerFuncSets []string
		handler        echo.HandlerFunc
	}

	handlerBundle struct {
		rootUrl        string
		apiVer         string
		handlerPathSets map[string]handlerEntity
		validate       *validator.Validate
	}

	handlerMiddle struct {}
)

var (
	DefaultHandlerBundle *handlerBundle
	localIpStr          string
)

// 初始化服务  绑定IP 端口
func init() {
	localIpStr = config.GetValue("system", "local_ip")

	DefaultHandlerBundle = &handlerBundle{
		rootUrl:        "",
		apiVer:         "",
		handlerPathSets: make(map[string]handlerEntity),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

}

func MainControl() {
	e := echo.New()

	e.Logger.SetLevel(echoLog.WARN)
	e.HTTPErrorHandler = HTTPErrorHandler

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	middle := &handlerMiddle{}
	e.Pre(process)
	e.Use(middle.trackId)
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection: "1; mode=block",
	}))

	for _, handler := range handlers.HandlerMappings {
		DefaultHandlerBundle.handlerPathSets[handler.AccessPath()] = handlerEntity{
			handlerFuncSets: handler.SupportMethod(),
			handler:        handler.Handler,
		}
	}

	for path, entity := range DefaultHandlerBundle.handlerPathSets {
		pathComponents := []string{
			DefaultHandlerBundle.rootUrl,
			path,
		}
		fullPath := strings.Join(pathComponents, "/")
		trimPath := strings.Trim(fullPath, "/")
		log.Info("MainControl", "support full path:", fullPath)
		log.Info("MainControl", "support trim path:", trimPath)
		e.Match(entity.handlerFuncSets, fullPath, entity.handler)
	}

	serverAddress := config.GetValue("system", "server_address")

	startServer(e, serverAddress)
}

// 开启服务
func startServer(e *echo.Echo, address string) {
	go func() {
		if err := e.Start(address); err != nil {
			log.Error("[startServer]server_stop", "err", err)
			os.Exit(1)
		}
	}()

	// after start log console
	time.Sleep(10 * time.Millisecond)
	e.Logger.SetOutput(log.GetLogger())

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	{{.lt}}-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("[startServer]server_stop", "err=", err)
	}
}

// middleware
func process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqBegin := time.Now()

		if err := next(c); err != nil {
			c.Error(err)
		}

		duration := time.Now().Sub(reqBegin)
		req := c.Request()

		statusCode := c.Response().Status
		// statsd metric
		if statusCode == http.StatusOK {
			//uri_split := strings.Split(req.RequestURI, "?")
			//path_split := strings.Split(uri_split[0], "/")
			//bucket := path_split[len(path_split)-1] + "_"
			//if path_split[len(path_split)-1] == "" {
			//	bucket = path_split[len(path_split)-2] + "_"
			//}
			//statsd.Statsd_metric(bucket, duration)
		}

		realIp := c.RealIP()
		remoteAddr := c.Request().RemoteAddr

		log.Info(
			"[process]stat",
			"method",
			req.Method,
			"url",
			req.RequestURI,
			"ip",
			realIp,
			"remote_addr",
			remoteAddr,
			"cost",
			duration.String(),
		)

		return nil
	}
}

func (p *handlerMiddle) trackId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		queryValues := req.URL.Query()
		var divide string = ""
		if "" == queryValues.Get("trackid") {
			reqUid := uuid.NewV4()
			if 0 != len(req.URL.RawQuery) {
				divide = "&"
			}
			req.URL.RawQuery = req.URL.RawQuery + divide + "trackid=" + reqUid.String()
		}
		return next(c)
	}
}

// copy DefaultHTTPErrorHandler to process the Security scanning of inner ip or UA eqs to OpenVAS
// DefaultHTTPErrorHandler is the default HTTP error handler. It sends a JSON response
// with status code.
func HTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // Issue #608
			if err := c.NoContent(code); err != nil {
				goto ERROR
			}
		} else {
			if err := c.JSON(code, msg); err != nil {
				goto ERROR
			}
		}
	}
	ERROR:
	isScanRequest, _ := regexp.Match("OpenVAS", []byte(c.Request().UserAgent()))
	if !isScanRequest {
		remoteIp := strings.Split(c.Request().RemoteAddr, ":")[0]
		if remoteIp == localIpStr {
			log.Error("[HTTPErrorHandler]", "err", err)
		} else if isInnerIp(remoteIp) {
			log.Info("[HTTPErrorHandler]", "err", err)
		} else {
			log.Error("[HTTPErrorHandler]", "err", err)
		}
	}
}

//10.0.0.0 167772160
//10.255.255.255 184549375
//172.16.0.0 2886729728
//172.31.255.255 2887778303
//192.168.0.0 3232235520
//192.168.255.255 3232301055
// not include local_ip
func isInnerIp(ipStr string) bool {
	ip := common.InetAton(ipStr)
	if (ip > 167772160 && ip {{.lt}} 184549375) ||
		(ip > 2886729728 && ip {{.lt}} 2887778303) ||
		(ip > 3232235520 && ip {{.lt}} 3232301055) {
		return true
	}
	return false
}
