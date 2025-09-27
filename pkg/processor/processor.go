package processor

import (
	parse "cut_utility/pkg/parsing"
	"fmt"
	"strings"
)

// ProcessLine - Извлечение поля из строки по разделителю
func ProcessLine(line string, config *parse.Config) {
	if config.Separated && !strings.Contains(line, config.Delimiter) {
		return
	}

	fields := strings.Split(line, config.Delimiter)
	var result []string

	for _, col := range config.Fields {
		if col >= len(fields) {
			return
		}
		result = append(result, fields[col])
	}

	if len(result) > 0 {
		fmt.Println(strings.Join(result, config.Delimiter))
	}
}
