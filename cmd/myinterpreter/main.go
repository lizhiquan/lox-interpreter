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

	if command != "tokenize" && command != "parse" && command != "evaluate" {
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

	switch command {
	case "tokenize":
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		for _, token := range tokens {
			fmt.Println(token)
		}

		if len(errs) > 0 {
			os.Exit(65)
		}

	case "parse":
		parser := lox.NewParser(tokens)
		expr, err := parser.Parse()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}

		printer := lox.AstPrinter{}
		fmt.Println(printer.Print(expr))

	case "evaluate":
		parser := lox.NewParser(tokens)
		expr, err := parser.Parse()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}

		var interpreter lox.Interpreter
		val, err := interpreter.Interpret(expr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(70)
		}

		fmt.Println(val)
	}
}
