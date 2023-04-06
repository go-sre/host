package middleware

import (
	"errors"
	"fmt"
	"github.com/go-sre/host/controller"
	"net/http"
)

func SignalHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ctrl controller.Controller
	if controller.IsEgressTraffic(r.URL.Query()) {
		ctrl = controller.EgressTable().LookupByName(r.URL.Query().Get(controller.RouteKey))
	} else {
		ctrl = controller.IngressTable().LookupByName(r.URL.Query().Get(controller.RouteKey))
	}
	if ctrl == nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New(fmt.Sprintf("invalid argument: route [%s] not found in [%s] table", r.URL.Query().Get(controller.RouteKey), r.URL.Query().Get(controller.TrafficKey)))
		w.Write([]byte(err.Error()))
		return
	}
	err := ctrl.Signal(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
