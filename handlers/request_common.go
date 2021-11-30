package handlers

import (
	"github.com/labstack/echo"
)

type (
	GeneratorHandler interface {
		AccessPath() string
		SupportMethod() []string
		Handler(cont echo.Context) error
	}

	CommonHandler struct {
		Path          string
		HandlerWorks map[string]func(echo.Context) error
	}
	BaseReqEntity struct {
		Trackid string `query:"trackid" validate:"lte_len=50"`
	}

	BaseResEntity struct {
		StatusCode string      `json:"code"`
		StatusDesc string      `json:"msg"`
		Data       interface{} `json:"data"`
	}
)

var (
	HandlerMappings = []GeneratorHandler{}
)

func (p *CommonHandler) AccessPath() string {
	return p.Path
}

func (p *CommonHandler) SupportMethod() (supportMethod []string) {
	supportMethod = []string{}
	for method := range p.HandlerWorks {
		supportMethod = append(supportMethod, method)
	}
	return
}

func (p *CommonHandler) Handler(cont echo.Context) error {
	method := cont.Request().Method
	err := p.HandlerWorks[method](cont)
	return err
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
