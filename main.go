package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	cmds := []Command{
		GetNewCommand(),
		GetSourceCommand(),
		GetResumeCommand(),
	}

	names := []string{}
	for _, cmd := range cmds {
		names = append(names, cmd.Name())
	}

	if len(os.Args) == 1 {
		log.Fatal(subCommandNotFoundError(names))
	}

	subCommand := os.Args[1]
	match := false
	for _, cmd := range cmds {
		if cmd.Name() == subCommand {
			if err := cmd.Parse(os.Args[2:]); err != nil {
				log.Fatal("Error: Failed to parse flags %v", err)
			}
			if err := cmd.Run(); err != nil {
				log.Fatal("Error: Failed to run %v", err)
			}
			match = true
		}
	}

	if !match {
		log.Fatal(subCommandNotFoundError(names))
	}
}

func subCommandNotFoundError(subcommands []string) error {
	return fmt.Errorf(
		"Error: No subcommand specified. Please provide one of the following subcommands: %s",
		strings.Join(subcommands, ", "),
	)
}
