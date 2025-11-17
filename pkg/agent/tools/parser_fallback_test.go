package tools

import (
	"encoding/xml"
	"testing"
)

func TestUnmarshalXMLWithFallback(t *testing.T) {
	type TestArgs struct {
		XMLName xml.Name `xml:"arguments"`
		Path    string   `xml:"path"`
		Pattern string   `xml:"pattern"`
		Text    string   `xml:"text"`
		URL     string   `xml:"url"`
		Content string   `xml:"content"`
		Outer   struct {
			Inner string `xml:"inner"`
		} `xml:"outer"`
	}

	tests := []struct {
		name      string
		input     string
		wantErr   bool
		checkFunc func(t *testing.T, result TestArgs)
	}{
		{
			name: "valid XML passes through unchanged",
			input: `<arguments>
				<path>src/main.go</path>
				<pattern>func.*main</pattern>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				if result.Path != "src/main.go" {
					t.Errorf("expected path=src/main.go, got %v", result.Path)
				}
				if result.Pattern != "func.*main" {
					t.Errorf("expected pattern=func.*main, got %v", result.Pattern)
				}
			},
		},
		{
			name: "unescaped single ampersand is fixed",
			input: `<arguments>
				<pattern>a & b</pattern>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				if result.Pattern != "a & b" {
					t.Errorf("expected pattern='a & b', got %v", result.Pattern)
				}
			},
		},
		{
			name: "unescaped double ampersand is fixed",
			input: `<arguments>
				<pattern>func.*&&.*return</pattern>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				if result.Pattern != "func.*&&.*return" {
					t.Errorf("expected pattern with &&, got %v", result.Pattern)
				}
			},
		},
		{
			name: "already escaped ampersand preserved",
			input: `<arguments>
				<text>a &amp; b</text>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				if result.Text != "a & b" {
					t.Errorf("expected text='a & b', got %v", result.Text)
				}
			},
		},
		{
			name: "multiple XML entities preserved",
			input: `<arguments>
				<text>&lt; &gt; &quot; &apos; &amp;</text>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				expected := "< > \" ' &"
				if result.Text != expected {
					t.Errorf("expected text=%q, got %v", expected, result.Text)
				}
			},
		},
		{
			name: "numeric entities preserved",
			input: `<arguments>
				<text>&#60; &#x3C; &#62; &#x3E;</text>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				expected := "< < > >"
				if result.Text != expected {
					t.Errorf("expected text=%q, got %v", expected, result.Text)
				}
			},
		},
		{
			name: "multiple unescaped ampersands in URL",
			input: `<arguments>
				<url>http://example.com?a=1&b=2&c=3</url>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				expected := "http://example.com?a=1&b=2&c=3"
				if result.URL != expected {
					t.Errorf("expected url=%q, got %v", expected, result.URL)
				}
			},
		},
		{
			name: "mixed escaped and unescaped ampersands",
			input: `<arguments>
				<text>Already &amp; escaped & not escaped</text>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				expected := "Already & escaped & not escaped"
				if result.Text != expected {
					t.Errorf("expected text=%q, got %v", expected, result.Text)
				}
			},
		},
		{
			name: "ampersand at start of value",
			input: `<arguments>
				<text>&value</text>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				if result.Text != "&value" {
					t.Errorf("expected text='&value', got %v", result.Text)
				}
			},
		},
		{
			name: "ampersand at end of value",
			input: `<arguments>
				<text>value&</text>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				if result.Text != "value&" {
					t.Errorf("expected text='value&', got %v", result.Text)
				}
			},
		},
		{
			name: "nested elements with unescaped ampersands",
			input: `<arguments>
				<outer>
					<inner>test & data</inner>
				</outer>
			</arguments>`,
			wantErr: false,
			checkFunc: func(t *testing.T, result TestArgs) {
				if result.Outer.Inner != "test & data" {
					t.Errorf("expected inner='test & data', got %v", result.Outer.Inner)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TestArgs
			err := UnmarshalXMLWithFallback([]byte(tt.input), &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXMLWithFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFunc != nil {
				tt.checkFunc(t, result)
			}
		})
	}
}

func TestEscapeUnescapedAmpersands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ampersands",
			input:    "<tag>hello world</tag>",
			expected: "<tag>hello world</tag>",
		},
		{
			name:     "single unescaped ampersand",
			input:    "<tag>a & b</tag>",
			expected: "<tag>a &amp; b</tag>",
		},
		{
			name:     "already escaped ampersand",
			input:    "<tag>a &amp; b</tag>",
			expected: "<tag>a &amp; b</tag>",
		},
		{
			name:     "double ampersand",
			input:    "<tag>a && b</tag>",
			expected: "<tag>a &amp;&amp; b</tag>",
		},
		{
			name:     "all standard entities",
			input:    "<tag>&amp; &lt; &gt; &quot; &apos;</tag>",
			expected: "<tag>&amp; &lt; &gt; &quot; &apos;</tag>",
		},
		{
			name:     "numeric character references",
			input:    "<tag>&#60; &#x3C;</tag>",
			expected: "<tag>&#60; &#x3C;</tag>",
		},
		{
			name:     "mixed escaped and unescaped",
			input:    "<tag>&amp; & &lt; & &gt;</tag>",
			expected: "<tag>&amp; &amp; &lt; &amp; &gt;</tag>",
		},
		{
			name:     "URL with query parameters",
			input:    "<url>http://example.com?a=1&b=2&c=3</url>",
			expected: "<url>http://example.com?a=1&amp;b=2&amp;c=3</url>",
		},
		{
			name:     "ampersand at boundaries",
			input:    "<tag>&start middle& &end</tag>",
			expected: "<tag>&amp;start middle&amp; &amp;end</tag>",
		},
		{
			name:     "complex regex pattern",
			input:    "<pattern>func.*&&.*return</pattern>",
			expected: "<pattern>func.*&amp;&amp;.*return</pattern>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeUnescapedAmpersands([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("escapeUnescapedAmpersands() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestUnmarshalXMLWithFallback_RealWorldScenarios(t *testing.T) {
	type SearchArgs struct {
		XMLName     xml.Name `xml:"arguments"`
		Path        string   `xml:"path"`
		Pattern     string   `xml:"pattern"`
		FilePattern string   `xml:"file_pattern"`
		URL         string   `xml:"url"`
		Content     string   `xml:"content"`
		Text        string   `xml:"text"`
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "search_files with regex containing ampersands",
			input: `<arguments>
				<path>src</path>
				<pattern>func.*&&.*return</pattern>
				<file_pattern>*.go</file_pattern>
			</arguments>`,
			wantErr: false,
		},
		{
			name: "URL with multiple query parameters",
			input: `<arguments>
				<url>https://api.example.com/v1/search?query=test&limit=10&offset=0</url>
			</arguments>`,
			wantErr: false,
		},
		{
			name: "code with logical AND operators",
			input: `<arguments>
				<content>if (condition1 && condition2 && condition3) {
	return true;
}</content>
			</arguments>`,
			wantErr: false,
		},
		{
			name: "mixed content with entities and unescaped",
			input: `<arguments>
				<text>Use &lt;tag&gt; for HTML & use && for AND</text>
			</arguments>`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result SearchArgs
			err := UnmarshalXMLWithFallback([]byte(tt.input), &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXMLWithFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
