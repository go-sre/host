package accessdata

import (
	"fmt"
)

func Example_IsDirectOperator() {
	op := Operator{Name: "test", Value: "   %"}
	fmt.Printf("test: IsDirectOperator() -> %v [value:%v]\n", IsDirectOperator(op), op.Value)

	op = Operator{Name: "test", Value: "%"}
	fmt.Printf("test: IsDirectOperator() -> %v [value:%v]\n", IsDirectOperator(op), op.Value)

	//Output:
	//test: IsDirectOperator() -> true [value:   %]
	//test: IsDirectOperator() -> false [value:%]
}

func Example_IsRequestOperator() {
	op := Operator{}
	ok := IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(<empty>) -> %v\n", ok)

	op = Operator{Name: " ", Value: " "}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(<empty>) -> %v\n", ok)

	op = Operator{Name: "", Value: "REQ "}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ(header"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ(header)"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ()"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ(1)%"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ(header-name)%"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	//Output:
	//test: IsRequestOperator(<empty>) -> false
	//test: IsRequestOperator(<empty>) -> false
	//test: IsRequestOperator({ REQ }) -> false
	//test: IsRequestOperator({ %REQ(header}) -> false
	//test: IsRequestOperator({ %REQ(header)}) -> false
	//test: IsRequestOperator({ %REQ()}) -> false
	//test: IsRequestOperator({ %REQ(1)%}) -> true
	//test: IsRequestOperator({ %REQ(header-name)%}) -> true

}

func Example_RequestOperatorHeaderName() {
	op := Operator{}
	name := RequestOperatorHeaderName(op)
	fmt.Printf("test: RequestOperatorHeaderName() -> %v [op:%v]\n", name, op.Value)

	op = Operator{Name: "", Value: "%REQ("}
	name = RequestOperatorHeaderName(op)
	fmt.Printf("test: RequestOperatorHeaderName() -> %v [op:%v]\n", name, op.Value)

	op = Operator{Name: "", Value: "%REQ()"}
	name = RequestOperatorHeaderName(op)
	fmt.Printf("test: RequestOperatorHeaderName() -> %v [op:%v]\n", name, op.Value)

	op = Operator{Name: "", Value: "%REQ()%"}
	name = RequestOperatorHeaderName(op)
	fmt.Printf("test: RequestOperatorHeaderName() -> %v [op:%v]\n", name, op.Value)

	op = Operator{Name: "", Value: "%REQ(1)%"}
	name = RequestOperatorHeaderName(op)
	fmt.Printf("test: RequestOperatorHeaderName() -> %v [op:%v]\n", name, op.Value)

	op = Operator{Name: "", Value: "%REQ(name)%"}
	name = RequestOperatorHeaderName(op)
	fmt.Printf("test: RequestOperatorHeaderName() -> %v [op:%v]\n", name, op.Value)

	//Output:
	//test: RequestOperatorHeaderName() ->  [op:]
	//test: RequestOperatorHeaderName() ->  [op:%REQ(]
	//test: RequestOperatorHeaderName() ->  [op:%REQ()]
	//test: RequestOperatorHeaderName() ->  [op:%REQ()%]
	//test: RequestOperatorHeaderName() -> 1 [op:%REQ(1)%]
	//test: RequestOperatorHeaderName() -> name [op:%REQ(name)%]

}

func Example_IsStringValue() {
	op := Operator{Name: "test", Value: "   %"}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: "%"}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: DurationOperator}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: TimeoutDurationOperator}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: RateLimitOperator}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: RateBurstOperator}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: ResponseStatusCodeOperator}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: ResponseBytesSentOperator}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	op = Operator{Name: "test", Value: ResponseBytesReceivedOperator}
	fmt.Printf("test: IsStringValue() -> %v [value:%v]\n", IsStringValue(op), op.Value)

	//Output:
	//test: IsStringValue() -> true [value:   %]
	//test: IsStringValue() -> true [value:%]
	//test: IsStringValue() -> false [value:%DURATION%]
	//test: IsStringValue() -> false [value:%TIMEOUT_DURATION%]
	//test: IsStringValue() -> false [value:%RATE_LIMIT%]
	//test: IsStringValue() -> false [value:%RATE_BURST%]
	//test: IsStringValue() -> false [value:%STATUS_CODE%]
	//test: IsStringValue() -> false [value:%BYTES_SENT%]
	//test: IsStringValue() -> false [value:%BYTES_RECEIVED%]

}
