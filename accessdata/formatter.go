package accessdata

// Formatter - template parameter for formatting
type Formatter interface {
	Format(items []Operator, data *Entry) string
}

type TextFormatter struct{}

func (TextFormatter) Format(items []Operator, data *Entry) string { return WriteText(items, data) }

type JsonFormatter struct{}

func (JsonFormatter) Format(items []Operator, data *Entry) string { return WriteJson(items, data) }
