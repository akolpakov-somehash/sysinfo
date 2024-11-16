package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Formatter interface {
	Format([]map[string]string) (string, error)
}

const (
	JSONFormat = "json"
	TextFormat = "text"
)

type JSONFormatter struct{}

func (j *JSONFormatter) Format(data []map[string]string) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

type TextFormatter struct{}

func (t *TextFormatter) Format(data []map[string]string) (string, error) {
	sb := strings.Builder{}
	for _, m := range data {
		title, ok := m["name"]
		hasTitle := 1
		if !ok {
			title = "Unknown"
			hasTitle = 0
		}
		sb.WriteString(fmt.Sprintf("Metrics for: %s\n", title))
		keys := make([]string, 0, len(m)-hasTitle)
		for k := range m {
			if k == "name" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for n := range keys {
			v := m[keys[n]]
			sb.WriteString(fmt.Sprintf("%s: %s\n", cases.Title(language.English).String(keys[n]), v))
		}
		sb.WriteString("\n----------------\n")
	}
	return sb.String(), nil
}

func NewFormatter(format string) Formatter {
	switch format {
	case JSONFormat:
		return &JSONFormatter{}
	case TextFormat:
		return &TextFormatter{}
	default:
		return &TextFormatter{}
	}
}
