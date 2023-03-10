package accessdata

import (
	"fmt"
	"strings"
)

const (
	markupNull   = "\"%v\":null"
	markupString = "\"%v\":\"%v\""
	markupValue  = "\"%v\":%v"
)

func WriteJson(items []Operator, data *Entry) string {
	if len(items) == 0 || data == nil {
		return "{}"
	}
	sb := strings.Builder{}
	for _, op := range items {
		writeMarkup(&sb, op.Name, data.Value(op.Value), IsStringValue(op))
	}
	sb.WriteString("}")
	return sb.String()
}

func writeMarkup(sb *strings.Builder, name, value string, stringValue bool) {
	if sb.Len() == 0 {
		sb.WriteString("{")
	} else {
		sb.WriteString(",")
	}
	if value == "" {
		sb.WriteString(fmt.Sprintf(markupNull, name))
	} else {
		format := markupString
		if !stringValue {
			format = markupValue
		}
		sb.WriteString(fmt.Sprintf(format, name, value))
	}
}

func WriteText(items []Operator, data *Entry) string {
	if len(items) == 0 || data == nil {
		return ""
	}
	sb := strings.Builder{}
	for i, op := range items {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(data.Value(op.Value))
	}

	return sb.String()
}
