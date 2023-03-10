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
