package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Sample   string `json:"sample,omitempty"`
}

type Schema struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) DetectFromJSON(data []byte) (*Schema, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	schema := &Schema{
		Name:   "detected_schema",
		Fields: make([]Field, 0, len(obj)),
	}

	for key, value := range obj {
		field := Field{
			Name:     key,
			Type:     d.inferType(value),
			Nullable: value == nil,
			Sample:   d.getSample(value),
		}
		schema.Fields = append(schema.Fields, field)
	}

	return schema, nil
}

func (d *Detector) inferType(value interface{}) string {
	if value == nil {
		return "string"
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		if d.isTimestamp(value.(string)) {
			return "timestamp"
		}
		return "string"
	case reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice:
		return "array"
	case reflect.Map:
		return "object"
	default:
		return "string"
	}
}

func (d *Detector) isTimestamp(s string) bool {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, format := range formats {
		if _, err := time.Parse(format, s); err == nil {
			return true
		}
	}
	return false
}

func (d *Detector) getSample(value interface{}) string {
	if value == nil {
		return ""
	}
	str := fmt.Sprintf("%v", value)
	if len(str) > 50 {
		return str[:47] + "..."
	}
	return str
}