package connector

import (
	"errors"
	"fmt"
	"github.com/go-sre/core/exchange"
	"github.com/go-sre/core/runtime"
	"net/http"
	url2 "net/url"
	"time"
)

var (
	pullLocInit      = pkgPath + "/initialize-pull"
	pullLocPull      = pkgPath + "/pull"
	pullErrorHandler runtime.ErrorHandleFn
	pullUrl          string
	pullClient       = http.DefaultClient
	done             chan bool
	request          *http.Request
)

func InitializePull[E runtime.ErrorHandler](uri string, newClient *http.Client) *runtime.Status {
	pullErrorHandler = runtime.NewErrorHandler[E]()
	if uri == "" {
		return pullErrorHandler(nil, pullLocInit, errors.New("invalid argument: uri is empty"))
	}
	u, err1 := url2.Parse(uri)
	if err1 != nil {
		return pullErrorHandler(nil, pullLocInit, err1)
	}
	pullUrl = u.String()
	var err error
	request, err = http.NewRequest("GET", pullUrl, nil)
	if err != nil {
		return pullErrorHandler(nil, pullLocInit, errors.New(fmt.Sprintf("invalid argument: upstream request error [%v]", err)))
	}
	if newClient != nil {
		pullClient = newClient
	}
	done = make(chan bool)
	pull(time.Second)
	return runtime.NewStatusOK()
}

func pull(d time.Duration) {
	tick := time.Tick(d)
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case <-tick:
			resp, err := pushClient.Do(request)
			if err != nil {
				pullErrorHandler(nil, pullLocPull, err)
			} else {
				_, err0 := exchange.ReadAll(resp.Body) //fmt.Println("Current time: ", t)
				if err0 != nil {
					pullErrorHandler(nil, pullLocPull, err0)
				}
			}
		default:
		}
	}
}


