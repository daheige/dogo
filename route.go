package dogo

import (
	"net/http"
	"strings"
)

type route struct {
	rules map[string]route
	h     handle
}

var Route = &route{rules: make(map[string]route)}

type handle func(c *Context)

func (r route) init(method string) {
	if _, ok := r.rules[method]; !ok {
		r.rules[method] = route{rules: make(map[string]route)}
	}
}

func (r route) Any(rule string, h handle) {
	r.init("*")
	(r.rules["*"]).addRoute(strings.Split(rule, "/"), h)
}

func (r route) Get(rule string, h handle) {
	r.init("GET")
	Dglog.Debug(h)
	(r.rules["GET"]).addRoute(strings.Split(rule, "/"), h)
}

func (r route) Post(rule string, h handle) {
	r.init("POST")
	(r.rules["POST"]).addRoute(strings.Split(rule, "/"), h)
}

func (r route) addRoute(rule []string, h handle) {
	Dglog.Debugf("add route %s ", rule)

	if len(rule) > 1 && rule[0] == "" {
		rule = rule[1:]
	} else {
		rule = rule
	}

	var nextRoute route
	if nextRouteTemp, ok := r.rules[rule[0]]; ok {
		nextRoute = nextRouteTemp
	} else {
		nextRoute = route{rules: make(map[string]route)}
	}

	if len(rule) == 1 {
		nextRoute.h = h
	} else {
		nextRoute.addRoute(rule[1:], h)
	}
	r.rules[rule[0]] = nextRoute
}

func (r route) checkRoute(request *http.Request) handle {
	method := request.Method
	url := strings.Split(request.RequestURI, "/")
	h := r.checkMethod(method, url)
	if h == nil {
		h = r.checkMethod("*", url)
	}
	return h
}

func (r route) checkMethod(method string, url []string) handle {
	return (r.rules[strings.ToUpper(method)]).checkUrl(url)
}

func (r route) checkUrl(url []string) handle {
	if len(url) <= 0 {
		return r.h
	} else if len(url) > 1 && url[0] == "" {
		url = url[1:]
	}
	return (r.rules[url[0]]).checkUrl(url[1:])
}
