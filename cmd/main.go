package main

import (
	"bufio"
	"fmt"
	"os"

	parse "cut_utility/pkg/parsing"
	"cut_utility/pkg/processor"
)

func main() {
	err := runApp()
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}

func runApp() error {
	config, err := parse.ParseArgs()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		processor.ProcessLine(line, config)
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
