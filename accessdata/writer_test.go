package accessdata

import (
	"fmt"
	"strings"
)

func Example_WriteJson() {
	sb := strings.Builder{}

	writeMarkup(&sb, "first", "string value", true)
	writeMarkup(&sb, "second", "100", false)
	writeMarkup(&sb, "third", "another string value", true)
	writeMarkup(&sb, "fourth", "true", false)
	writeMarkup(&sb, "null-value", "", false)
	sb.WriteString("}")

	fmt.Printf("test: writeMarkup() -> [%v]\n", sb.String())

	//Output:
	//test: writeMarkup() -> [{"first":"string value","second":100,"third":"another string value","fourth":true,"null-value":null}]
}
