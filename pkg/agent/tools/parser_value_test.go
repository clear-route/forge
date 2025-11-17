package tools

import (
	"math"
	"testing"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		typeStr  string
	}{
		// Booleans
		{"bool true lowercase", "true", true, "bool"},
		{"bool false lowercase", "false", false, "bool"},
		{"bool true uppercase", "TRUE", true, "bool"},
		{"bool false uppercase", "FALSE", false, "bool"},
		{"bool true mixed case", "True", true, "bool"},
		{"bool false mixed case", "False", false, "bool"},
		{"bool true with spaces", "  true  ", true, "bool"},
		{"bool false with spaces", "  false  ", false, "bool"},

		// Null
		{"null lowercase", "null", nil, "nil"},
		{"null uppercase", "NULL", nil, "nil"},
		{"null mixed case", "Null", nil, "nil"},
		{"null with spaces", "  null  ", nil, "nil"},

		// Integers
		{"int zero", "0", 0, "int"},
		{"int positive", "123", 123, "int"},
		{"int negative", "-42", -42, "int"},
		{"int with spaces", "  456  ", 456, "int"},
		{"int large positive", "2147483647", 2147483647, "int"},
		{"int large negative", "-2147483648", -2147483648, "int"},
		// Note: On 64-bit systems, int64 max/min values fit in int
		{"int max value", "9223372036854775807", 9223372036854775807, "int"},
		{"int min value", "-9223372036854775808", -9223372036854775808, "int"},

		// Floats
		{"float simple", "123.45", 123.45, "float64"},
		{"float negative", "-3.14", -3.14, "float64"},
		{"float zero", "0.0", 0.0, "float64"},
		{"float scientific", "1.23e10", 1.23e10, "float64"},
		{"float scientific negative", "-2.5e-3", -2.5e-3, "float64"},
		{"float with spaces", "  99.99  ", 99.99, "float64"},

		// Strings (anything that doesn't match above patterns)
		{"string simple", "hello", "hello", "string"},
		{"string empty", "", "", "string"},
		{"string with spaces", "  hello world  ", "hello world", "string"},
		{"string alphanumeric", "abc123", "abc123", "string"},
		{"string number-like", "123abc", "123abc", "string"},
		{"string bool-like", "truthy", "truthy", "string"},
		{"string null-like", "nullable", "nullable", "string"},
		{"string quoted bool", `"true"`, `"true"`, "string"},
		{"string quoted number", `"123"`, `"123"`, "string"},
		{"string with special chars", "hello@world.com", "hello@world.com", "string"},
		{"string path", "/path/to/file.txt", "/path/to/file.txt", "string"},

		// Edge cases
		{"leading zeros", "007", 7, "int"},
		{"negative zero", "-0", 0, "int"},
		{"float no decimal", "42.", 42.0, "float64"},
		// Note: ParseFloat successfully parses these special values
		{"infinity positive", "+Inf", math.Inf(1), "float64"},
		{"infinity negative", "-Inf", math.Inf(-1), "float64"},
		{"not a number", "NaN", math.NaN(), "float64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseValue(tt.input)

			// Type assertion based on expected type
			switch tt.typeStr {
			case "bool":
				v, ok := result.(bool)
				if !ok {
					t.Errorf("parseValue(%q) returned type %T, want bool", tt.input, result)
					return
				}
				if v != tt.expected.(bool) {
					t.Errorf("parseValue(%q) = %v, want %v", tt.input, v, tt.expected)
				}

			case "int":
				v, ok := result.(int)
				if !ok {
					t.Errorf("parseValue(%q) returned type %T, want int", tt.input, result)
					return
				}
				if v != tt.expected.(int) {
					t.Errorf("parseValue(%q) = %v, want %v", tt.input, v, tt.expected)
				}

			case "int64":
				v, ok := result.(int64)
				if !ok {
					t.Errorf("parseValue(%q) returned type %T, want int64", tt.input, result)
					return
				}
				if v != tt.expected.(int64) {
					t.Errorf("parseValue(%q) = %v, want %v", tt.input, v, tt.expected)
				}

			case "float64":
				v, ok := result.(float64)
				if !ok {
					t.Errorf("parseValue(%q) returned type %T, want float64", tt.input, result)
					return
				}
				expected := tt.expected.(float64)
				// Special handling for NaN
				if math.IsNaN(expected) {
					if !math.IsNaN(v) {
						t.Errorf("parseValue(%q) = %v, want NaN", tt.input, v)
					}
				} else if math.IsInf(expected, 0) {
					// Special handling for infinity
					if !math.IsInf(v, 0) || math.Signbit(v) != math.Signbit(expected) {
						t.Errorf("parseValue(%q) = %v, want %v", tt.input, v, expected)
					}
				} else {
					// Use approximate comparison for regular floats
					if math.Abs(v-expected) > 1e-10 {
						t.Errorf("parseValue(%q) = %v, want %v", tt.input, v, expected)
					}
				}

			case "string":
				v, ok := result.(string)
				if !ok {
					t.Errorf("parseValue(%q) returned type %T, want string", tt.input, result)
					return
				}
				if v != tt.expected.(string) {
					t.Errorf("parseValue(%q) = %q, want %q", tt.input, v, tt.expected)
				}

			case "nil":
				if result != nil {
					t.Errorf("parseValue(%q) = %v, want nil", tt.input, result)
				}

			default:
				t.Fatalf("Unknown type string: %s", tt.typeStr)
			}
		})
	}
}

func TestParseValueNestedXML(t *testing.T) {
	// Test that nested XML is parsed recursively
	xmlContent := "<name>John</name><age>30</age>"
	result := parseValue(xmlContent)

	nested, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	if nested["name"] != "John" {
		t.Errorf("Expected name='John', got %v", nested["name"])
	}

	if nested["age"] != 30 {
		t.Errorf("Expected age=30 (int), got %v (%T)", nested["age"], nested["age"])
	}
}

func TestParseValueTypePreservation(t *testing.T) {
	// Verify that type conversion happens correctly for common cases
	tests := []struct {
		input        string
		expectedType string
	}{
		{"true", "bool"},
		{"false", "bool"},
		{"null", "nil"},
		{"123", "int"},
		{"123.45", "float64"},
		{"hello", "string"},
	}

	for _, tt := range tests {
		result := parseValue(tt.input)
		var actualType string

		switch result.(type) {
		case bool:
			actualType = "bool"
		case int:
			actualType = "int"
		case int64:
			actualType = "int64"
		case float64:
			actualType = "float64"
		case string:
			actualType = "string"
		case nil:
			actualType = "nil"
		default:
			actualType = "unknown"
		}

		if actualType != tt.expectedType {
			t.Errorf("parseValue(%q) returned type %s, want %s", tt.input, actualType, tt.expectedType)
		}
	}
}