package prompts

import (
	"strings"
	"testing"
)

func TestGenerateXMLExample(t *testing.T) {
	t.Run("SimpleStringParameter", func(t *testing.T) {
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "A simple message",
				},
			},
			"required": []string{"message"},
		}

		result := GenerateXMLExample(schema, "test_tool")

		if !strings.Contains(result, "<tool_name>test_tool</tool_name>") {
			t.Error("Expected tool_name in result")
		}
		if !strings.Contains(result, "<message>") {
			t.Error("Expected message parameter in result")
		}
	})

	t.Run("CDATAForContentFields", func(t *testing.T) {
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"content": map[string]interface{}{
					"type":        "string",
					"description": "File content",
				},
			},
			"required": []string{"content"},
		}

		result := GenerateXMLExample(schema, "write_file")

		if !strings.Contains(result, "<![CDATA[") {
			t.Error("Expected CDATA for content field")
		}
	})

	t.Run("NestedArrayOfObjects", func(t *testing.T) {
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"edits": map[string]interface{}{
					"type":        "array",
					"description": "List of edits",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"search": map[string]interface{}{
								"type":        "string",
								"description": "Search text",
							},
							"replace": map[string]interface{}{
								"type":        "string",
								"description": "Replace text",
							},
						},
					},
				},
			},
			"required": []string{"edits"},
		}

		result := GenerateXMLExample(schema, "apply_diff")

		// Should have nested structure
		if !strings.Contains(result, "<edits>") {
			t.Error("Expected edits array")
		}
		if !strings.Contains(result, "<edit>") {
			t.Error("Expected edit element (singular)")
		}
		if !strings.Contains(result, "<search>") {
			t.Error("Expected search element")
		}
		if !strings.Contains(result, "<replace>") {
			t.Error("Expected replace element")
		}

		// Should NOT have CDATA wrapping the entire edits structure
		if strings.Contains(result, "<edits><![CDATA[") {
			t.Error("Should NOT wrap array structure in CDATA")
		}
	})

	t.Run("BooleanParameter", func(t *testing.T) {
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"recursive": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to recurse",
				},
			},
			"required": []string{"recursive"},
		}

		result := GenerateXMLExample(schema, "list_files")

		if !strings.Contains(result, "<recursive>true</recursive>") {
			t.Error("Expected boolean true value")
		}
	})

	t.Run("NumericParameters", func(t *testing.T) {
		schema := map[string]interface{}{
			"properties": map[string]interface{}{
				"count": map[string]interface{}{
					"type":        "integer",
					"description": "Count value",
				},
				"ratio": map[string]interface{}{
					"type":        "number",
					"description": "Ratio value",
				},
			},
			"required": []string{"count", "ratio"},
		}

		result := GenerateXMLExample(schema, "test_tool")

		if !strings.Contains(result, "<count>42</count>") {
			t.Error("Expected integer example")
		}
		if !strings.Contains(result, "<ratio>3.14</ratio>") {
			t.Error("Expected float example")
		}
	})
}
