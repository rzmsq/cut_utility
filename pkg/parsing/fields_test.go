package parsing

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestParseFieldsRaw(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
		hasError bool
	}{
		{
			name:     "single field",
			input:    "1",
			expected: []int{0},
			hasError: false,
		},
		{
			name:     "multiple fields",
			input:    "1,3,5",
			expected: []int{0, 2, 4},
			hasError: false,
		},
		{
			name:     "range of fields",
			input:    "1-3",
			expected: []int{0, 1, 2},
			hasError: false,
		},
		{
			name:     "mixed fields and ranges",
			input:    "1,3-5,7",
			expected: []int{0, 2, 3, 4, 6},
			hasError: false,
		},
		{
			name:     "duplicate fields",
			input:    "1,1,2",
			expected: []int{0, 1},
			hasError: false,
		},
		{
			name:     "unordered fields",
			input:    "3,1,2",
			expected: []int{0, 1, 2},
			hasError: false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
			hasError: false,
		},
		{
			name:     "spaces in input",
			input:    " 1 , 2 , 3 ",
			expected: []int{0, 1, 2},
			hasError: false,
		},
		{
			name:     "invalid field - zero",
			input:    "0",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid field - negative",
			input:    "-1",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid field - non-numeric",
			input:    "abc",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid range - too many parts",
			input:    "1-2-3",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid range - non-numeric start",
			input:    "a-3",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid range - non-numeric end",
			input:    "1-b",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid range - zero start",
			input:    "0-3",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid range - zero end",
			input:    "1-0",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFieldsRaw(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For input %q, expected %v, but got %v", tt.input, tt.expected, result)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Config
		hasError bool
	}{
		{
			name: "basic configuration",
			args: []string{"program", "-f", "1,2,3", "-d", ",", "-s"},
			expected: &Config{
				Fields:    []int{0, 1, 2},
				Delimiter: ",",
				Separated: true,
			},
			hasError: false,
		},
		{
			name: "default delimiter",
			args: []string{"program", "-f", "1,2"},
			expected: &Config{
				Fields:    []int{0, 1},
				Delimiter: "\t",
				Separated: false,
			},
			hasError: false,
		},
		{
			name: "range fields",
			args: []string{"program", "-f", "1-3"},
			expected: &Config{
				Fields:    []int{0, 1, 2},
				Delimiter: "\t",
				Separated: false,
			},
			hasError: false,
		},
		{
			name:     "invalid fields",
			args:     []string{"program", "-f", "0"},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальные аргументы командной строки
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Сбрасываем флаги перед каждым тестом
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Устанавливаем новые аргументы
			os.Args = tt.args

			result, err := ParseArgs()

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for args %v, but got none", tt.args)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for args %v: %v", tt.args, err)
				return
			}

			if !reflect.DeepEqual(result.Fields, tt.expected.Fields) {
				t.Errorf("For args %v, expected fields %v, but got %v", tt.args, tt.expected.Fields, result.Fields)
			}

			if result.Delimiter != tt.expected.Delimiter {
				t.Errorf("For args %v, expected delimiter %q, but got %q", tt.args, tt.expected.Delimiter, result.Delimiter)
			}

			if result.Separated != tt.expected.Separated {
				t.Errorf("For args %v, expected separated %v, but got %v", tt.args, tt.expected.Separated, result.Separated)
			}
		})
	}
}
