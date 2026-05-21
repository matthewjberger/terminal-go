// Command zork is a small data-oriented text adventure.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"terminal-go/command"
	"terminal-go/parse"
	"terminal-go/world"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(in io.Reader, out io.Writer) error {
	w := world.NewDemo()

	fmt.Fprintln(out, "terminal-go zork. Type 'help' for verbs, 'save' / 'load' to persist, 'quit' to leave.")
	command.DescribeRoom(w, w.PlayerRoom, out)

	scanner := bufio.NewScanner(in)
	for !w.Quit && !w.Won {
		fmt.Fprint(out, "\n> ")
		if !scanner.Scan() {
			fmt.Fprintln(out)
			return scanner.Err()
		}
		command.Run(w, parse.Tokenize(scanner.Text()), out)
	}
	return scanner.Err()
}
