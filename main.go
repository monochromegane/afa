package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
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

	defaultSubCommandIdx := 0
	for i, cmd := range cmds {
		if cmd.Default() {
			defaultSubCommandIdx = i
		}
	}
	defaultSubCommand := []string{cmds[defaultSubCommandIdx].Name(), "-script"}

	args := os.Args[1:]
	if len(args) == 0 {
		args = defaultSubCommand
	}

	flagSet := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	flagSet.Usage = func() {}
	flagSetOutput := flagSet.Output()
	flagSet.SetOutput(io.Discard)
	ver := flagSet.Bool("version", false, "Display version")
	if err := flagSet.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(flagSetOutput, "Usage of %s:\n", cmdName)
			flagSet.PrintDefaults()
			for _, cmd := range cmds {
				fmt.Fprintf(flagSetOutput, "  %s\n\t%s\n", cmd.Name(), cmd.Description())
			}
			os.Exit(0)
		}
		// No subcommand and unknown flag
		args = append(defaultSubCommand, args...)
	}
	if *ver {
		fmt.Fprintf(flagSetOutput, "%s v%s (rev:%s)\n", cmdName, version, revision)
		os.Exit(0)
	}

	names := []string{}
	for _, cmd := range cmds {
		names = append(names, cmd.Name())
	}

	if len(os.Args) == 1 {
		log.Fatal(subCommandNotFoundError(names))
	}

	subCommand := args[0]
	match := false
	for _, cmd := range cmds {
		if cmd.Name() == subCommand {
			if err := cmd.Parse(args[1:]); err != nil {
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
