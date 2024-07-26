package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/lox"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := lox.NewScanner(string(fileContents))
	tokens, errs := scanner.ScanTokens()
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}
	for _, token := range tokens {
		fmt.Println(token)
	}

	if len(errs) > 0 {
		os.Exit(65)
	}
}
