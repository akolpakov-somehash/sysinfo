package formatter

import (
	"fmt"
	"strings"
)

type Formatter interface {
	Format([]map[string]string) string
}

type JSONFormatter struct{}

func (j *JSONFormatter) Format(data []map[string]string) string {
	return ""
}

type TextFormatter struct{}

func (t *TextFormatter) Format(data []map[string]string) string {
	sb := strings.Builder{}
	for _, m := range data {
		for k, v := range m {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
		sb.WriteString("\n----------------\n")
	}
	return sb.String()
}

func NewFormatter(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "text":
		return &TextFormatter{}
	default:
		return &TextFormatter{}
	}
}
