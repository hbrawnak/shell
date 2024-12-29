package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	for {
		printPrompt()
		command := readCommand() // Wait for user input
		if command == "" {
			continue // Ignore empty commands
		}
		execute(command)
	}
}

func printPrompt() {
	fmt.Fprint(os.Stdout, "$ ")
}

func readCommand() string {
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		handleError(err)
	}
	return strings.TrimSpace(command)
}

func execute(input string) {
	parts := parseCommand(input)
	command := parts[0]
	args := parts[1:]

	switch command {
	case "exit":
		exit(args)
	case "echo":
		echo(args)
	case "type":
		types(args)
	case "pwd":
		pwd()
	case "cd":
		cd(args)
	default:
		runExternalCommand(command, args)
	}
}

func parseCommand(input string) []string {
	// Regex to match single-quoted strings or unquoted parts
	re := regexp.MustCompile(`'[^']*'|[^' \t]+`)
	matches := re.FindAllString(input, -1)

	for i, match := range matches {
		if len(match) > 1 && match[0] == '\'' && match[len(match)-1] == '\'' {
			matches[i] = match[1 : len(match)-1]
		}
	}
	return matches
}

func handleError(err error) {
	if err.Error() == "EOF" {
		fmt.Fprintln(os.Stdout, "exit")
		return
	}

	fmt.Fprintln(os.Stderr, "Error reading input:", err)
	os.Exit(1)
}
