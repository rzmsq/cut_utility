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
	for i, col := range config.Fields {
		if col >= len(fields) {
			return
		}
		printLine(fields[col], i, config)
	}
	return
}

// printLine - Вывод строки
func printLine(field string, indx int, config *parse.Config) {
	if indx >= len(field)-1 {
		fmt.Print(field)
	} else {
		fmt.Print(field + config.Delimiter)
	}
	fmt.Println()
}
