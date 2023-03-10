package accesslog

import (
	"fmt"
	"github.com/gotemplates/host/accessdata"
)

func ExampleOutputHandler() {
	fmt.Printf("test: Output[NilOutputHandler,data.TextFormatter](nil,nil)\n")
	logTest[NilOutputHandler, accessdata.TextFormatter](nil, nil)

	fmt.Printf("test: Output[DebugOutputHandler,data.JsonFormatter](operators,data)\n")
	ops := []accessdata.Operator{{"error", "message"}}
	logTest[DebugOutputHandler, accessdata.JsonFormatter](ops, accessdata.NewEmptyEntry())

	fmt.Printf("test: Output[TestOutputHandler,data.JsonFormatter](nil,nil)\n")
	logTest[TestOutputHandler, accessdata.JsonFormatter](nil, nil)

	fmt.Printf("test: Output[TestOutputHandler,data.JsonFormatter](ops,data)\n")
	logTest[TestOutputHandler, accessdata.JsonFormatter](ops, accessdata.NewEmptyEntry())

	fmt.Printf("test: Output[LogOutputHandler,data.JsonFormatter](ops,data)\n")
	logTest[LogOutputHandler, accessdata.JsonFormatter](ops, accessdata.NewEmptyEntry())

	//Output:
	//test: Output[NilOutputHandler,data.TextFormatter](nil,nil)
	//test: Output[DebugOutputHandler,data.JsonFormatter](operators,data)
	//{"error":"message"}
	//test: Output[TestOutputHandler,data.JsonFormatter](nil,nil)
	//test: Write() -> [{}]
	//test: Output[TestOutputHandler,data.JsonFormatter](ops,data)
	//test: Write() -> [{"error":"message"}]
	//test: Output[LogOutputHandler,data.JsonFormatter](ops,data)

}

func logTest[O OutputHandler, F accessdata.Formatter](items []accessdata.Operator, data *accessdata.Entry) {
	var o O
	var f F
	o.Write(items, data, f)
}
