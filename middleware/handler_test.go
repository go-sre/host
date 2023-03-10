package middleware

import (
	"fmt"
	"net/http"
)

func ExampleTimeoutHandler() {

}

func _ExampleMiddleware() {
	m := http.NewServeMux()
	if m != nil {
	}
	//m.Handler()
	fmt.Printf("test () -> [%v]\n", "results")

	//Output:
	//fail
}
