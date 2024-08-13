package main

import (
	"flag"
)

type Command interface {
	Name() string
	Parse([]string) error
	Run() error
}

type NewCommand struct {
	flagSet  *flag.FlagSet
	aiForAll *AIForAll
}

func (c NewCommand) Name() string { return "new" }

func (c *NewCommand) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *NewCommand) Run() error {
	return c.aiForAll.New()
}

type SourceCommand struct {
	flagSet  *flag.FlagSet
	aiForAll *AIForAll
}

func (c SourceCommand) Name() string { return "source" }

func (c *SourceCommand) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *SourceCommand) Run() error {
	return c.aiForAll.Source()
}

type ResumeCommand struct {
	flagSet  *flag.FlagSet
	aiForAll *AIForAll
}

func (c ResumeCommand) Name() string { return "resume" }

func (c *ResumeCommand) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *ResumeCommand) Run() error {
	return c.aiForAll.Resume()
}

func GetNewCommand() Command {
	flagSet := flag.NewFlagSet("new", flag.ExitOnError)
	aiForAll := &AIForAll{}

	setBasicFlags(aiForAll, flagSet)

	return &NewCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}
}

func GetSourceCommand() Command {
	flagSet := flag.NewFlagSet("source", flag.ExitOnError)
	aiForAll := &AIForAll{}

	setBasicFlags(aiForAll, flagSet)

	return &SourceCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}
}

func GetResumeCommand() Command {
	flagSet := flag.NewFlagSet("resume", flag.ExitOnError)
	aiForAll := &AIForAll{}

	setBasicFlags(aiForAll, flagSet)

	return &ResumeCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}
}

func setBasicFlags(aiForAll *AIForAll, flagSet *flag.FlagSet) {
	flagSet.StringVar(&aiForAll.Project, "p", "default", "Name of project.")
}
