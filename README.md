## Custom Shell Implementation in Go

This project is a custom shell built in Go, capable of handling various built-in and external commands. It demonstrates core shell functionalities like command parsing, handling quoted strings, and executing both internal and external commands.

## Features

- **Built-in Commands**:
    - `echo`: Prints arguments separated by spaces.
    - `cd`: Changes the current working directory.
    - `pwd`: Displays the current working directory.
    - `type`: Identifies whether a command is built-in or an external executable.
    - `exit`: Exits the shell.

- **External Command Execution**:
    - Executes commands available in the system's `PATH`.

- **Command Parsing**:
    - Supports handling of single and double-quoted strings.
    - Escapes special characters and spaces appropriately.

## How It Works

The shell reads commands from the user, parses them into individual arguments, and executes them based on their type (built-in or external). The project is modular, with functions dedicated to handling specific commands and utilities.

### Example Usage

```bash
$ echo "Hello, world!"
Hello, world!

$ pwd
/home/user

$ cd /tmp
$ pwd
/tmp

$ type echo
echo is a shell builtin

$ type ls
ls is /bin/ls

$ exit
```

## Author
Developed by [Md Habibur Rahman](https://habib.im).