package processor

import (
	"bytes"
	"cut_utility/pkg/parsing"
	"io"
	"os"
	"testing"
)

// captureOutput захватывает вывод в stdout для тестирования
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return ""
	}
	return buf.String()
}

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		config         *parsing.Config
		expectedOutput string
	}{
		{
			name: "single field with tab delimiter",
			line: "field1\tfield2\tfield3",
			config: &parsing.Config{
				Fields:    []int{0},
				Delimiter: "\t",
				Separated: false,
			},
			expectedOutput: "field1\n",
		},
		{
			name: "multiple fields with tab delimiter",
			line: "field1\tfield2\tfield3",
			config: &parsing.Config{
				Fields:    []int{0, 2},
				Delimiter: "\t",
				Separated: false,
			},
			expectedOutput: "field1\tfield3\n",
		},
		{
			name: "single field with comma delimiter",
			line: "field1,field2,field3",
			config: &parsing.Config{
				Fields:    []int{1},
				Delimiter: ",",
				Separated: false,
			},
			expectedOutput: "field2\n",
		},
		{
			name: "multiple fields with comma delimiter",
			line: "field1,field2,field3,field4",
			config: &parsing.Config{
				Fields:    []int{0, 2, 3},
				Delimiter: ",",
				Separated: false,
			},
			expectedOutput: "field1,field3,field4\n",
		},
		{
			name: "field index out of bounds",
			line: "field1\tfield2",
			config: &parsing.Config{
				Fields:    []int{0, 5}, // field 5 doesn't exist
				Delimiter: "\t",
				Separated: false,
			},
			expectedOutput: "field1\n", // should return nothing when field is out of bounds
		},
		{
			name: "separated mode - line contains delimiter",
			line: "field1,field2,field3",
			config: &parsing.Config{
				Fields:    []int{0, 1},
				Delimiter: ",",
				Separated: true,
			},
			expectedOutput: "field1,field2\n",
		},
		{
			name: "separated mode - line doesn't contain delimiter",
			line: "field1 field2 field3",
			config: &parsing.Config{
				Fields:    []int{0, 1},
				Delimiter: ",",
				Separated: true,
			},
			expectedOutput: "", // should return nothing when separator is not found
		},
		{
			name: "empty fields",
			line: "field1,,field3",
			config: &parsing.Config{
				Fields:    []int{0, 1, 2},
				Delimiter: ",",
				Separated: false,
			},
			expectedOutput: "field1,,field3\n",
		},
		{
			name: "single character fields",
			line: "a\tb\tc\td",
			config: &parsing.Config{
				Fields:    []int{1, 3},
				Delimiter: "\t",
				Separated: false,
			},
			expectedOutput: "b\td\n",
		},
		{
			name: "duplicate field indices",
			line: "field1,field2,field3",
			config: &parsing.Config{
				Fields:    []int{0, 0, 1}, // duplicate field 0
				Delimiter: ",",
				Separated: false,
			},
			expectedOutput: "field1,field1,field2\n",
		},
		{
			name: "empty line",
			line: "",
			config: &parsing.Config{
				Fields:    []int{0},
				Delimiter: "\t",
				Separated: false,
			},
			expectedOutput: "\n", // empty string split gives one empty element
		},
		{
			name: "line with only delimiter",
			line: "\t\t",
			config: &parsing.Config{
				Fields:    []int{0, 1, 2},
				Delimiter: "\t",
				Separated: false,
			},
			expectedOutput: "\t\t\n", // three empty fields
		},
		{
			name: "custom delimiter",
			line: "field1|field2|field3",
			config: &parsing.Config{
				Fields:    []int{0, 2},
				Delimiter: "|",
				Separated: false,
			},
			expectedOutput: "field1|field3\n",
		},
		{
			name: "whitespace delimiter",
			line: "field1 field2 field3",
			config: &parsing.Config{
				Fields:    []int{1},
				Delimiter: " ",
				Separated: false,
			},
			expectedOutput: "field2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				ProcessLine(tt.line, tt.config)
			})

			if output != tt.expectedOutput {
				t.Errorf("For line %q with config %+v, expected output %q, but got %q",
					tt.line, tt.config, tt.expectedOutput, output)
			}
		})
	}
}

func TestProcessLineEdgeCases(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		// Тест на панику при передаче nil конфига
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when config is nil, but didn't get one")
			}
		}()

		ProcessLine("test", nil)
	})

	t.Run("empty fields slice", func(t *testing.T) {
		config := &parsing.Config{
			Fields:    []int{},
			Delimiter: "\t",
			Separated: false,
		}

		output := captureOutput(func() {
			ProcessLine("field1\tfield2", config)
		})

		if output != "" {
			t.Errorf("Expected empty output for empty fields slice, but got %q", output)
		}
	})

	t.Run("very large field index", func(t *testing.T) {
		config := &parsing.Config{
			Fields:    []int{1000},
			Delimiter: "\t",
			Separated: false,
		}

		output := captureOutput(func() {
			ProcessLine("field1\tfield2", config)
		})

		if output != "" {
			t.Errorf("Expected empty output for very large field index, but got %q", output)
		}
	})

	t.Run("multi-character delimiter", func(t *testing.T) {
		config := &parsing.Config{
			Fields:    []int{0, 1},
			Delimiter: "||",
			Separated: false,
		}

		output := captureOutput(func() {
			ProcessLine("field1||field2||field3", config)
		})

		expected := "field1||field2\n"
		if output != expected {
			t.Errorf("Expected output %q for multi-character delimiter, but got %q", expected, output)
		}
	})
}
