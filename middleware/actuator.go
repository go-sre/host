package middleware

import (
	"errors"
	"fmt"
	"github.com/go-sre/host/controller"
	"net/http"
	"net/url"
	"strings"
)

func ActuatorHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.URL == nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New(fmt.Sprintf("invalid argument: request URL is nil"))
		w.Write([]byte(err.Error()))
		return
	}
	if r.URL.Query() == nil || len(r.URL.Query()) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New(fmt.Sprintf("invalid argument: request URL does not contain any query arguments"))
		w.Write([]byte(err.Error()))
		return
	}
	traffic, route, behavior, err := parseUrl(r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var ctrl controller.Controller
	if traffic == controller.EgressTraffic {
		ctrl = controller.EgressTable().LookupByName(route)
	} else {
		ctrl = controller.IngressTable().LookupByName(route)
	}
	if ctrl == nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.New(fmt.Sprintf("invalid argument: route [%s] not found in [%s] table", route, traffic))
		w.Write([]byte(err.Error()))
		return
	}
	values := r.URL.Query()
	values.Add(controller.BehaviorKey, behavior)
	err = ctrl.Signal(values)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func parseUrl(url *url.URL) (traffic string, route string, behavior string, err error) {
	if url == nil {
		return "", "", "", errors.New("invalid argument: request URL is nil")
	}
	tokens := strings.Split(url.Path, "/")
	if len(tokens) == 1 {
		return "", "", "", errors.New("invalid argument: request URL path is empty")
	}
	index := 0
	if tokens[0] == "" {
		index++
	}
	if tokens[index] != "actuator" {
		return "", "", "", errors.New("invalid argument: request URL path does not start with 'actuator'")
	}
	if len(tokens) <= 2 {
		return "", "", "", errors.New("invalid argument: request URL path does not contain traffic type")
	}

	index++
	traffic = tokens[index]
	if traffic != controller.IngressTraffic && traffic != controller.EgressTraffic {
		return traffic, "", "", errors.New(fmt.Sprintf("invalid argument: request URL path does not contain valid traffic type [%v]", traffic))
	}

	if len(tokens) <= 3 {
		return traffic, "", "", errors.New("invalid argument: request URL path does not contain route name")
	}
	index++
	route = tokens[index]
	if len(tokens) <= 4 {
		return traffic, route, "", errors.New("invalid argument: request URL path does not contain behavior name")
	}
	index++
	behavior = tokens[index]
	return traffic, route, behavior, nil
}
