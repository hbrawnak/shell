package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
		handleExit(args)
	case "echo":
		handleEcho(args)
	case "type":
		handleType(args)
	case "pwd":
		handlePwd()
	case "cd":
		handleCd(args)
	default:
		runExternalCommand(command, args)
	}
}

func parseCommand(input string) []string {
	re := regexp.MustCompile(`'[^']*'|"([^"\\]*(\\.[^"\\]*)*)"|\\.|[^'"\t ]+`)
	matches := re.FindAllString(input, -1)

	for i, match := range matches {
		switch {
		case strings.HasPrefix(match, `'`) && strings.HasSuffix(match, `'`):
			matches[i] = match[1 : len(match)-1]
		case strings.HasPrefix(match, `"`) && strings.HasSuffix(match, `"`):
			unescaped := match[1 : len(match)-1]
			unescaped = strings.ReplaceAll(unescaped, `\'`, `'`)
			unescaped = strings.ReplaceAll(unescaped, `\"`, `"`)
			unescaped = strings.ReplaceAll(unescaped, `\\`, `\`)
			matches[i] = unescaped
		case strings.HasPrefix(match, `\`) || strings.HasSuffix(match, `\`):
			unescaped := strings.ReplaceAll(match, `\"`, `"`)
			unescaped = strings.ReplaceAll(unescaped, `\'`, `'`)
			unescaped = strings.ReplaceAll(unescaped, `'\`, `'`)
			unescaped = strings.ReplaceAll(unescaped, `"\`, `"`)
			unescaped = strings.ReplaceAll(unescaped, `\`, ``)
			matches[i] = strings.ReplaceAll(unescaped, ` `, "")
		default:
			// Replace escaped spaces with a literal space
			// unescaped := strings.TrimPrefix(match, `\`)
			// matches[i] = strings.TrimSuffix(unescaped, `\`)
			matches[i] = match
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

// handleExit handles the exit command. It exits the program if the argument is "0", otherwise,
// it prints an error message.
func handleExit(args []string) {
	if len(args) == 1 && args[0] == "0" {
		os.Exit(0)
	}
	fmt.Println("Invalid exit command format")
}

// handleEcho prints the arguments joined by a space.
func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// handleType prints information about a command: whether it's a shell builtin or an executable found
// in the system PATH.
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

// handlePwd prints the current working directory.
func handlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting current directory:", err)
		return
	}
	fmt.Println(dir)
}

// handleCd changes the current working directory. It handles relative paths and the home directory symbol ('~').
func handleCd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout, "cd: missing argument")
		return
	}

	dir := args[0]

	// Handle '~' as the home directory
	if dir == "~" {
		homeDir, exists := os.LookupEnv("HOME")
		if !exists || homeDir == "" {
			fmt.Fprintln(os.Stdout, "cd: $HOME not set")
			return
		}
		dir = homeDir
	}

	cleanPath := path.Clean(dir)

	if !path.IsAbs(cleanPath) {
		dir, _ = os.Getwd()
		cleanPath = path.Join(dir, cleanPath)
	}

	if err := os.Chdir(cleanPath); err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", cleanPath)
	}
}

// findExecutablePath searches for the executable in the directories listed in
// the PATH environment variable.
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

// runExternalCommand runs an external command and prints its output or error
// to the standard output/error.
func runExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("%s: command not found\n", command)
	}
}
