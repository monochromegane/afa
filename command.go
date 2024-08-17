package main

import (
	"flag"
	"io"
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
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	c.aiForAll.Files = c.flagSet.Args()
	return nil
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
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	c.aiForAll.Files = c.flagSet.Args()
	return nil
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
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	c.aiForAll.Files = c.flagSet.Args()
	return nil
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

	if err := setBasicFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}

	flagSet.StringVar(&aiForAll.SystemPromptTemplate, "S", "default", "Name of system prompt template.")
	flagSet.StringVar(&aiForAll.Model, "M", "gpt-4o-mini", "Name of Model.")

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

	if err := setBasicFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}

	flagSet.StringVar(&aiForAll.SessionName, "l", "", "Log name of session.")

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

	if err := setBasicFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}

	return &ResumeCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func setBasicFlags(aiForAll *AIForAll, flagSet *flag.FlagSet) error {
	flagSet.StringVar(&aiForAll.Message, "m", "", "Message as initial prompt.")
	flagSet.StringVar(&aiForAll.UserPromptTemplate, "U", "default", "Name of user prompt template.")
	flagSet.BoolVar(&aiForAll.Interactive, "i", false, "Runs in interactive mode; set to false when standard input is passed.")
	flagSet.BoolVar(&aiForAll.Stream, "s", false, "Runs in stream mode.")
	flagSet.Func("R", "Resume based on the identifier of latest session. (default \"$PPID\")", func(runsOn string) error {
		aiForAll.RunsOn = runsOn
		return nil
	})

	if hasStdin() {
		inputStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		aiForAll.MessageStdin = string(inputStdin)
		aiForAll.Interactive = false
	}

	return nil
}

func hasStdin() bool {
	if stat, err := os.Stdin.Stat(); err == nil {
		return (stat.Mode() & os.ModeCharDevice) == 0
	}
	return false
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
