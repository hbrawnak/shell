package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	case "pwd":
		handlePwd()
	case "cd":
		cd(args)
	default:
		runExternalCommand(command, args)
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
	case "echo", "type", "exit", "pwd":
		fmt.Printf("%s is a shell builtin\n", args[0])
	default:
		if path := findExecutablePath(args[0]); path != "" {
			fmt.Printf("%s is %s\n", args[0], path)
		} else {
			fmt.Printf("%s: not found\n", args[0])
		}
	}
}

func findExecutablePath(command string) string {
	pathEnv := os.Getenv("PATH")
	directories := strings.Split(pathEnv, ":")

	for _, dir := range directories {
		fullPath := filepath.Join(dir, command)
		if fileInfo, err := os.Stat(fullPath); err == nil {
			if fileInfo.Mode().IsRegular() && (fileInfo.Mode().Perm()&0111 != 0) {
				return fullPath
			}
		}
	}

	return ""
}

func handlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting current directory:", err)
		return
	}
	fmt.Println(dir)
}

func cd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout, "cd: missing argument")
		return
	}

	dir := args[0]
	cleanPath := path.Clean(dir)

	if !path.IsAbs(cleanPath) {
		dir, _ = os.Getwd()
		cleanPath = path.Join(dir, cleanPath)
	}

	if err := os.Chdir(cleanPath); err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", cleanPath)
	}
}

func runExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("%s: command not found\n", command)
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
