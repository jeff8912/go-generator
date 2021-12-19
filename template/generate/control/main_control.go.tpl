package control

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/context"

	"{{.packagePrefix}}/config"
	"{{.packagePrefix}}/handler"
    "{{.packagePrefix}}/log"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
)

type (
	handleEntity struct {
		handleFuncSets []string
		handler        echo.HandlerFunc
	}

	handleBundle struct {
		rootUrl        string
		handlePathSets map[string]handleEntity
		validate       *validator.Validate
	}
)

var (
	DefaultHandleBundle *handleBundle
)

// 初始化服务  绑定IP 端口
func init() {
	DefaultHandleBundle = &handleBundle{
		rootUrl:        "",
		handlePathSets: make(map[string]handleEntity),
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

	e.Pre(process)
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection: "1; mode=block",
	}))

	for _, h := range handler.HandlerMappings {
		DefaultHandleBundle.handlePathSets[h.AccessPath()] = handleEntity{
			handleFuncSets: h.SupportMethod(),
			handler:        h.Handler,
		}
	}

	for path, entity := range DefaultHandleBundle.handlePathSets {
		pathComponents := []string{
			DefaultHandleBundle.rootUrl,
			path,
		}
		fullPath := strings.Join(pathComponents, "/")
		log.Info("main_control", "support full path:", fullPath)
		e.Match(entity.handleFuncSets, fullPath, entity.handler)
	}

	serverAddress := config.GetValue("system", "server_address")

	startServer(e, serverAddress)
}

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

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	{{ .lt }}-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("main_control", "event", "server_stop", "err", err)
	}
}

// process middleware
func process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		beginTime := time.Now()

		if err := next(c); err != nil {
			c.Error(err)
		}

		req := c.Request()

		log.Info(
			"[go-generator][main_control]",
			"method",
			req.Method,
			"url",
			req.RequestURI,
			"ip",
			c.RealIP(),
			"remote_addr",
			c.Request().RemoteAddr,
			"cost",
			time.Now().Sub(beginTime).String(),
		)

		return nil
	}
}

// HTTPErrorHandler copy DefaultHTTPErrorHandler to process the Security scanning of inner ip or UA eqs to OpenVAS
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
		if c.Request().Method == echo.HEAD {
			if err := c.NoContent(code); err != nil {
				log.Error("main_control[HTTPErrorHandler]", "err", err)
			}
		} else {
			if err := c.JSON(code, msg); err != nil {
				log.Error("main_control[HTTPErrorHandler]", "err", err)
			}
		}
	}
}
