package parsing

import (
	"errors"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	errInvalidFields = errors.New("invalid fields")
)

// Config содержит конфигурацию для парсинга полей из строк
type Config struct {
	Fields    []int
	Delimiter string
	Separated bool
}

// FieldSet представляет множество номеров полей для предотвращения дублирования
type FieldSet map[int]struct{}

// ParseArgs парсит аргументы командной строки и возвращает конфигурацию
func ParseArgs() (*Config, error) {
	var config Config
	var fieldRaw string

	flag.StringVar(&fieldRaw, "f", "", "Numbers of columns")
	flag.StringVar(&config.Delimiter, "d", "\t", "Delimiter of columns")
	flag.BoolVar(&config.Separated, "s", false, "Separated columns")
	flag.Parse()

	var err error
	config.Fields, err = parseFieldsRaw(fieldRaw)
	if err != nil {
		return nil, fmt.Errorf("[cut error] -f [fields]: %w", err)
	}

	return &config, nil
}

// parseFieldsRaw преобразует строку с номерами полей в слайс интов
func parseFieldsRaw(fieldsRaw string) ([]int, error) {
	var fields []int
	fieldSet := make(FieldSet)
	parts := strings.Split(fieldsRaw, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeOfField := strings.Split(part, "-")
			if len(rangeOfField) != 2 {
				return nil, errInvalidFields
			}
			start, err := strconv.Atoi(rangeOfField[0])
			if err != nil || start < 1 {
				return nil, errInvalidFields
			}
			end, err := strconv.Atoi(rangeOfField[1])
			if err != nil || end < 1 {
				return nil, errInvalidFields
			}
			if start > end {
				return nil, errInvalidFields
			}
			for i := start; i <= end; i++ {
				fieldSet[i-1] = struct{}{}
			}
		} else {
			n, err := strconv.Atoi(part)
			if err != nil || n < 1 {
				return nil, errInvalidFields
			}
			fieldSet[n-1] = struct{}{}
		}
	}

	for key := range fieldSet {
		fields = append(fields, key)
	}

	sort.Ints(fields)

	return fields, nil
}
