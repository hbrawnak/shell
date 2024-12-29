package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func exit(args []string) {
	if len(args) == 1 && args[0] == "0" {
		os.Exit(0)
	}
	fmt.Println("Invalid exit command format")
}

func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func types(args []string) {
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

func pwd() {
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

func runExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("%s: command not found\n", command)
	}
}
