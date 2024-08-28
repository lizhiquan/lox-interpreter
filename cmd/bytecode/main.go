package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"interpreter/lox"
)

func main() {
	argc := len(os.Args)

	if argc == 1 {
		repl()
	} else if argc == 2 {
		runFile(os.Args[1])
	} else {
		fmt.Fprintln(os.Stderr, "Usage: ./lox [path]")
		os.Exit(64)
	}
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		if scanner.Scan() {
			line := scanner.Text()
			_ = lox.Interpret(line)
		} else {
			if errors.Is(scanner.Err(), io.EOF) {
				break
			}
			panic(scanner.Err())
		}
	}
}

func runFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file: %s\n", err)
		os.Exit(74)
	}
	defer file.Close()

	source, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err)
		os.Exit(74)
	}

	err = lox.Interpret(string(source))
	if errors.Is(err, lox.ErrInterpretCompile) {
		os.Exit(65)
	}
	if errors.Is(err, lox.ErrInterpretRuntime) {
		os.Exit(70)
	}
}
