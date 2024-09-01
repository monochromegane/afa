package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	initCommand, err := GetInitCommand()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Failed to get init command. %v", err))
	}
	newCommand, err := GetNewCommand()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Failed to get new command. %v", err))
	}
	sourceCommand, err := GetSourceCommand()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Failed to get source command. %v", err))
	}
	resumeCommand, err := GetResumeCommand()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Failed to get resume command. %v", err))
	}
	listCommand, err := GetListCommand()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Failed to get list command. %v", err))
	}
	showCommand, err := GetShowCommand()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Failed to get show command. %v", err))
	}

	cmds := []Command{
		initCommand,
		newCommand,
		sourceCommand,
		resumeCommand,
		listCommand,
		showCommand,
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
				log.Fatal(fmt.Sprintf("Error: Failed to parse flags. %v", err))
			}
			if err := cmd.Run(); err != nil {
				log.Fatal(fmt.Sprintf("Error: Failed to run. %v", err))
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
		"Error: No subcommand specified. Please provide one of the following subcommands: %s.",
		strings.Join(subcommands, ", "),
	)
}
