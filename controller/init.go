package controller

import (
	"encoding/json"
	"errors"
)

const (
	DefaultIngressRouteName = "default-ingress"
	DefaultEgressRouteName  = "default-egress"
)

// ReadRoutes - read routes from the []byte representation of a route configuration
func ReadRoutes(buf []byte) ([]Route, error) {
	var config []RouteConfig

	if buf == nil {
		return nil, errors.New("invalid argument: buffer is nil")
	}
	err1 := json.Unmarshal(buf, &config)
	if err1 != nil {
		return nil, err1
	}
	var routes []Route
	for _, c := range config {
		r, err := NewRouteFromConfig(c)
		if err != nil {
			return nil, err
		}
		routes = append(routes, r)
	}
	return routes, nil
}

// AddEgressRoutes - read the routes from the []byte and create the EgressTable controller entries
func AddEgressRoutes(buf []byte) ([]Route, []error) {
	routes, err := ReadRoutes(buf)
	if err != nil {
		return routes, []error{err}
	}
	var errs []error
	for _, r := range routes {
		switch r.Name {
		case DefaultEgressRouteName:
			errs = EgressTable.SetDefaultController(r)
		default:
			errs = EgressTable.AddController(r)
		}
		if len(errs) > 0 {
			return nil, errs
		}
	}
	return routes, nil
}

// AddIngressRoutes - read the routes from the []byte and create the IngressTable controller entries
func AddIngressRoutes(buf []byte) ([]Route, []error) {
	routes, err := ReadRoutes(buf)
	if err != nil {
		return nil, []error{err}
	}
	var errs []error
	for _, r := range routes {
		switch r.Name {
		case HostControllerName:
			errs = IngressTable.SetHostController(r)
		case DefaultIngressRouteName:
			errs = IngressTable.SetDefaultController(r)
		default:
			errs = IngressTable.AddController(r)
		}
		if len(errs) > 0 {
			return nil, errs
		}
	}
	return routes, nil
}
