package main

import (
	"flag"
	"os"
	"path"
)

type Command interface {
	Name() string
	Parse([]string) error
	Run() error
}

type InitCommand struct {
	flagSet  *flag.FlagSet
	aiForAll *AIForAll
}

func (c InitCommand) Name() string { return "init" }

func (c *InitCommand) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *InitCommand) Run() error {
	return c.aiForAll.Init()
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

func GetInitCommand() (Command, error) {
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	return &InitCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetNewCommand() (Command, error) {
	flagSet := flag.NewFlagSet("new", flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	setBasicFlags(aiForAll, flagSet)

	return &NewCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetSourceCommand() (Command, error) {
	flagSet := flag.NewFlagSet("source", flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	setBasicFlags(aiForAll, flagSet)

	return &SourceCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetResumeCommand() (Command, error) {
	flagSet := flag.NewFlagSet("resume", flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	setBasicFlags(aiForAll, flagSet)

	return &ResumeCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func setBasicFlags(aiForAll *AIForAll, flagSet *flag.FlagSet) {
}

func newAIForAll() (*AIForAll, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	return NewAIForAll(
		path.Join(configDir, "afa"),
		path.Join(cacheDir, "afa"),
	), nil
}
