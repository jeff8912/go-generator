package control

import (
	"go-generator/handlers"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/context"

	"go-generator/common"
	"go-generator/config"
	"go-generator/log"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echoLog "github.com/labstack/gommon/log"
)

var (
	DefaultHandleBundle *handleBundle
	localIpStr          string
)

// 初始化服务  绑定IP 端口
func init() {
	localIpStr = config.GetValue("system", "local_ip")

	DefaultHandleBundle = &handleBundle{
		rootUrl:        "",
		apiVer:         "",
		handlePathSets: make(map[string]handleEntity),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

}

type (
	handleEntity struct {
		handleFuncSets []string
		handler        echo.HandlerFunc
	}

	handleBundle struct {
		rootUrl        string
		apiVer         string
		handlePathSets map[string]handleEntity
		validate       *validator.Validate
	}
)
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func MainControl() {
	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/html/*.html")),
	}
	e.Renderer = renderer

	e.Logger.SetLevel(echoLog.WARN)
	e.HTTPErrorHandler = HTTPErrorHandler

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	e.Pre(Process)
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection: "1; mode=block",
	}))

	for _, handler := range handlers.HandlerMappings {
		DefaultHandleBundle.handlePathSets[handler.AccessPath()] = handleEntity{
			handleFuncSets: handler.SupportMethod(),
			handler:        handler.Handler,
		}
	}

	for path, entity := range DefaultHandleBundle.handlePathSets {
		pathComponents := []string{
			DefaultHandleBundle.rootUrl,
			path,
		}
		fullPath := strings.Join(pathComponents, "/")
		trimPath := strings.Trim(fullPath, "/")
		log.Info("main_control", "support full path:", fullPath)
		log.Info("main_control", "support trim path:", trimPath)
		e.Match(entity.handleFuncSets, fullPath, entity.handler)
	}

	serverAddress := config.GetValue("system", "server_address")

	startServer(e, serverAddress)
}

// 开启服务
func startServer(e *echo.Echo, address string) {
	go func() {
		if err := e.Start(address); err != nil {
			log.Error("main_control", "event", "server_stop", "err", err)
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
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("main_control", "event", "server_stop", "err", err)
	}
}

// middleware
func Process(next echo.HandlerFunc) echo.HandlerFunc {
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
			"[go-generator][main_control]",
			"event",
			"stat",
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
	if (ip > 167772160 && ip < 184549375) ||
		(ip > 2886729728 && ip < 2887778303) ||
		(ip > 3232235520 && ip < 3232301055) {
		return true
	}
	return false
}
