package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			handleError(err)
		}

		command = strings.TrimSpace(command)
		if command == "" {
			continue // Ignore empty commands
		}

		executeCommand(command)
	}
}

func executeCommand(input string) {
	parts := strings.Fields(input)
	command := parts[0]
	args := parts[1:]

	switch command {
	case "exit":
		handleExit(args)
	case "echo":
		handleEcho(args)
	case "type":
		handleType(args)
	default:
		fmt.Printf("%s: command not found\n", command)
	}
}

func handleExit(args []string) {
	if len(args) == 1 && args[0] == "0" {
		os.Exit(0)
	}
	fmt.Println("Invalid exit command format")
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handleType(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: type <command>")
		return
	}

	switch args[0] {
	case "echo", "type", "exit":
		fmt.Printf("%s is a shell builtin\n", args[0])
	default:
		fmt.Printf("%s: not found\n", args[0])
	}
}

func handleError(err error) {
	if err.Error() == "EOF" {
		fmt.Fprintln(os.Stdout, "exit")
		return
	}

	fmt.Fprintln(os.Stderr, "Error reading input:", err)
	os.Exit(1)
}
