package accessdata

import (
	"fmt"
	"time"
)

func _ExampleFmtTimestamp() {
	t := time.Now().UTC()
	s := FmtTimestamp(t)
	fmt.Printf("test: FmtTimestamp() -> [%v]\n", s)

	t2, err := ParseTimestamp(s)
	fmt.Printf("test: ParseTimestamp() -> [%v] [%v]\n", FmtTimestamp(t2), err)

	//Output:

}
