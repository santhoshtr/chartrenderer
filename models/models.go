package models

import "fmt"

type ChartType string

const (
	Line    ChartType = "line"
	Bar     ChartType = "bar"
	Scatter ChartType = "scatter"
	Pie     ChartType = "pie"
	Area    ChartType = "area"
)

func (t ChartType) String() string {
	return string(t)
}

// MarshalJSON implements the json.Marshaler interface for the ChartType type.
// It converts the ChartType value to a JSON-encoded string.
func (t ChartType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *ChartType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"line"`:
		*t = Line
	case `"bar"`:
		*t = Bar
	case `"area"`:
		*t = Area
	case `"scatter"`:
		*t = Scatter
	case `"pie"`:
		*t = Pie
	default:
		return fmt.Errorf("invalid ChartType value: %s", string(data))
	}
	return nil
}

type FieldDataType string

const (
	String FieldDataType = "string"
	Number FieldDataType = "number"
)

// String returns the string representation of the FieldDataType.
// It converts the FieldDataType to its corresponding string value.
func (t FieldDataType) String() string {
	return string(t)
}

// MarshalJSON implements the json.Marshaler interface for the FieldDataType type.
// It converts the FieldDataType value to a JSON-encoded string.
func (t FieldDataType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *FieldDataType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"string"`:
		*t = String
	case `"number"`:
		*t = Number
	default:
		return fmt.Errorf("invalid FieldDataType value: %s", string(data))
	}
	return nil
}

type Axis struct {
	Title       map[string]string `json:"title"`
	Angle       int               `json:"angle"`
	Type        string            `json:"type"`
	ShowSymbols bool              `json:"showSymbols"`
	Source      string            `json:"source"`
}

type ChartDefinition struct {
	Version     int               `json:"version"`
	License     License           `json:"license"`
	Source      string            `json:"source"`
	Type        ChartType         `json:"type"`
	Title       map[string]string `json:"title"`
	XAxis       Axis              `json:"xAxis"`
	YAxis       Axis              `json:"yAxis"`
	Interpolate string            `json:"interpolate"`
}

type Field struct {
	Name  string            `json:"name"`
	Type  FieldDataType     `json:"type"`
	Title map[string]string `json:"title"`
}

type Schema struct {
	Fields []Field `json:"fields"`
}

type License struct {
	Code string `json:"code"`
	Text string `json:"text"`
	URL  string `json:"url"`
}

type ChartData struct {
	License     License           `json:"license"`
	Description map[string]string `json:"description"`
	Schema      Schema            `json:"schema"`
	Data        [][]interface{}   `json:"data"`
}
