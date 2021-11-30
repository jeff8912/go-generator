package handlers

import (
    "{{.packagePrefix}}/log"
    "github.com/go-playground/validator"

    "github.com/labstack/echo"
    "time"
    "{{.packagePrefix}}/constant"
)

type (
    CommonHandler struct {
        Path         string
        HandlerWorks map[string]func(echo.Context) error
    }

    Handler interface {
        AccessPath() string
        SupportMethod() []string
        Handler(cont echo.Context) error
    }
)

var(
    HandlerMappings = []Handler{}
    validate = validator.New()
)

func (p *CommonHandler) AccessPath() string {
    return p.Path
}

func (p *CommonHandler) SupportMethod() []string {
    supMethods := []string{}
    for method := range p.HandlerWorks {
        supMethods = append(supMethods, method)
    }
    return supMethods
}

func (p *CommonHandler) Handler(cont echo.Context) error {
    method := cont.Request().Method
    result := p.HandlerWorks[method](cont)
    return result
}

func (p *CommonHandler) requestMapping(path, method string, handlerWork func(echo.Context) error) {
    p.Path = path
    if p.HandlerWorks == nil {
        p.HandlerWorks = make(map[string]func(echo.Context) error)
        HandlerMappings = append(HandlerMappings, p)
    }
    p.HandlerWorks[method] = handlerWork
}

func (p *CommonHandler) getMapping(path string, handlerWork func(echo.Context) error) {
    p.requestMapping(path, echo.GET, handlerWork)
}

func (p *CommonHandler) postMapping(path string, handlerWork func(echo.Context) error) {
    p.requestMapping(path, echo.POST, handlerWork)
}

func (p *CommonHandler) putMapping(path string, handlerWork func(echo.Context) error) {
    p.requestMapping(path, echo.PUT, handlerWork)
}

func (p *CommonHandler) deleteMapping(path string, handlerWork func(echo.Context) error) {
    p.requestMapping(path, echo.DELETE, handlerWork)
}

// 读取请求体
func (p *CommonHandler) readBody(context echo.Context, body interface{}) (errCode int) {
    beginTime := time.Now()
    err := context.Bind(body)
    if err != nil {
        log.Error("[readBody]Bind fail", "queryString", context.QueryString(), "body", body, "err", err)
        return constant.BIND_PARAM_ERROR
    }

    if err := validate.Struct(body); err != nil {
        log.Error("[readBody]Validate fail", "queryString", context.QueryString(), "body", body, "err", err)
        return constant.VALIDATE_PARAM_ERROR
    }

    if log.GetLogLevel() == "debug" {
        log.Debug("[readBody]success", "cost", time.Now().Sub(beginTime).String())
    }
    return constant.SUCCESS
}

func getSessionId(cont echo.Context) string {
    sessionId, _ := cont.Request().Cookie("sessionid")
    return sessionId.Value
}