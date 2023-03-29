package controller

import "fmt"

func Example_ParseState() {
	s := "  timeout : 35, statusCode : 504 "
	names, values := ParseState(s)
	fmt.Printf("test: ParseState() -> [names:%v] values:%v\n", names, values)

	//Output:
	//test: ParseState() -> [names:[timeout statusCode]] values:[35 504]

}

func Example_ExtractState() {
	name := "invalid"
	state := "  timeout : 35, statusCode : 504 "
	value := ExtractState(state, name)
	fmt.Printf("test: ExtractState(%v) -> [%v]\n", name, value)

	name = "timeout"
	value = ExtractState(state, name)
	fmt.Printf("test: ExtractState(%v) -> [%v]\n", name, value)

	name = "statusCode"
	value = ExtractState(state, name)
	fmt.Printf("test: ExtractState(%v) -> [%v]\n", name, value)

	//Output:
	//test: ExtractState(invalid) -> []
	//test: ExtractState(timeout) -> [35]
	//test: ExtractState(statusCode) -> [504]

}

func ExampleConvertDuration() {
	s := ""
	duration, err := ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "  "
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "12as"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1000"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1000s"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1000m"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "1m"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	s = "10ms"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	//t := time.Microsecond * 100
	//fmt.Printf("test: time.String %v\n", t.String())

	s = "10µs"
	duration, err = ConvertDuration(s)
	fmt.Printf("test: ConvertDuration(\"%v\") [err:%v] [duration:%v]\n", s, err, duration)

	//Output:
	//test: ConvertDuration("") [err:<nil>] [duration:0s]
	//test: ConvertDuration("  ") [err:strconv.Atoi: parsing "  ": invalid syntax] [duration:0s]
	//test: ConvertDuration("12as") [err:strconv.Atoi: parsing "12a": invalid syntax] [duration:0s]
	//test: ConvertDuration("1000") [err:<nil>] [duration:16m40s]
	//test: ConvertDuration("1000s") [err:<nil>] [duration:16m40s]
	//test: ConvertDuration("1000m") [err:<nil>] [duration:16h40m0s]
	//test: ConvertDuration("1m") [err:<nil>] [duration:1m0s]
	//test: ConvertDuration("10ms") [err:<nil>] [duration:10ms]
	//test: ConvertDuration("10µs") [err:<nil>] [duration:10µs]

}
