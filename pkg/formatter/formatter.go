package formatter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sysinfo/internal/metrics"

	"github.com/dustin/go-humanize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Formatter is an interface that defines a method to format metrics data.
type Formatter interface {
	Format([][]metrics.Metric) (string, error)
}

const (
	// JSONFormat represents the JSON format type.
	JSONFormat = "json"
	// TextFormat represents the text format type.
	TextFormat = "text"
)

// JSONFormatter is a struct that implements the Formatter interface for JSON format.
type JSONFormatter struct{}

// Format formats the metrics data to JSON format.
func (j *JSONFormatter) Format(data [][]metrics.Metric) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	return string(result), nil
}

// TextFormatter is a struct that implements the Formatter interface for text format.
type TextFormatter struct {
	capitalSplitRegex *regexp.Regexp
	titleCaser        cases.Caser
}

// Format formats the given metrics data into a human-readable text string.
func (t *TextFormatter) Format(data [][]metrics.Metric) (string, error) {
	resultBuilder := strings.Builder{}
	sectionBuilder := strings.Builder{}
	for _, metricsList := range data {
		for _, m := range metricsList {
			if m.Type == metrics.TypeTitle {
				tempStr := sectionBuilder.String()
				sectionBuilder.Reset()
				sectionBuilder.WriteString(fmt.Sprintf("Metrics for %s\n-------\n", m.Value))
				sectionBuilder.WriteString(tempStr)
				continue
			}
			v := m.Value
			if m.Type == metrics.TypeByte {
				vUint, err := strconv.ParseUint(m.Value, 10, 64)
				if err != nil {
					return "", fmt.Errorf("failed to parse byte value '%s': %w", m.Value, err)
				}
				v = humanize.Bytes(vUint)
			}
			if m.Type == metrics.TypePer {
				vFloat, err := strconv.ParseFloat(m.Value, 64)
				if err != nil {
					return "", fmt.Errorf("failed to parse percentage value '%s': %w", m.Value, err)
				}
				v = fmt.Sprintf("%.2f%%", vFloat)
			}
			fmt.Fprintf(&sectionBuilder, "%s: %s\n", t.splitAndTitleCase(m.Name), v)
		}
		sectionBuilder.WriteString("\n")
		resultBuilder.WriteString(sectionBuilder.String())
		sectionBuilder.Reset()
	}
	return resultBuilder.String(), nil
}

// splitAndTitleCase splits a string by capitalized letters.
// For example, "modelName" will be split into "model Name".
func (t *TextFormatter) splitAndTitleCase(input string) string {
	output := t.capitalSplitRegex.ReplaceAllString(input, "$1 $2")
	return t.titleCaser.String(strings.ToLower(output))
}

var capitalSplitRegex = regexp.MustCompile(`([a-z])([A-Z])`)

// NewFormatter creates a new Formatter based on the given format.
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case JSONFormat:
		return &JSONFormatter{}, nil
	case TextFormat:
		return &TextFormatter{
			capitalSplitRegex: capitalSplitRegex,
			titleCaser:        cases.Title(language.English),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported format '%s', valid formats are 'json' or 'text'", format)
	}
}
