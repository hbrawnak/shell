package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		// Uncomment this block to pass the first stage
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {

			if err.Error() == "EOF" {
				fmt.Fprintln(os.Stdout, "exit")
				return
			}

			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command = strings.TrimSpace(command)

		parseCommand(command)
	}
}

func parseCommand(command string) {
	parts := strings.Fields(command) // Split the command into parts
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "exit":
		if len(parts) > 1 && parts[1] == "0" {
			os.Exit(0)
		} else {
			fmt.Println("Invalid exit command format")
		}
	case "echo":
		fmt.Println(strings.Join(parts[1:], " "))
	case "type":
		parseType(parts[1])
	default:
		if command != "" {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func parseType(typeStr string) {
	switch typeStr {
	case "echo", "type", "exit":
		fmt.Printf("%s is a shell builtin\n", typeStr)
	default:
		fmt.Printf("%s: not found\n", typeStr)
	}
}
