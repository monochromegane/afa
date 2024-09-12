package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const cmdName = "afa"

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

	flagSet := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Usage of %s:\n", cmdName)
		flagSet.PrintDefaults()
		for _, cmd := range cmds {
			fmt.Fprintf(flagSet.Output(), "  %s\n\t%s\n", cmd.Name(), cmd.Description())
		}
	}
	ver := flagSet.Bool("version", false, "Display version")
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		os.Exit(2)
	}
	if *ver {
		fmt.Fprintf(flagSet.Output(), "%s v%s (rev:%s)\n", cmdName, version, revision)
		os.Exit(0)
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
