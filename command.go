package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
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
	if c.aiForAll.WorkSpace.IsNotExist() {
		return workSpaceNotExistError()
	}
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
	if c.aiForAll.WorkSpace.IsNotExist() {
		return workSpaceNotExistError()
	}
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
	if c.aiForAll.WorkSpace.IsNotExist() {
		return workSpaceNotExistError()
	}
	return c.aiForAll.Resume()
}

type ListCommand struct {
	flagSet  *flag.FlagSet
	aiForAll *AIForAll
}

func (c ListCommand) Name() string { return "list" }

func (c *ListCommand) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *ListCommand) Run() error {
	return c.aiForAll.List()
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

	if err := setBasicChatFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}

	flagSet.StringVar(
		&aiForAll.Option.Chat.SystemPromptTemplate,
		"s",
		aiForAll.Option.Chat.SystemPromptTemplate,
		"Name of system prompt template.",
	)
	flagSet.StringVar(
		&aiForAll.Option.Chat.Model,
		"m",
		aiForAll.Option.Chat.Model,
		"Name of Model.",
	)
	flagSet.StringVar(
		&aiForAll.Option.Chat.Schema,
		"j",
		aiForAll.Option.Chat.Schema,
		"Name of JSON schema for response format.",
	)

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

	if err := setBasicChatFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}

	flagSet.StringVar(
		&aiForAll.SessionName,
		"l",
		aiForAll.SessionName,
		"Log name of session.",
	)
	flagSet.BoolVar(
		&aiForAll.Option.Chat.WithHistory,
		"H",
		aiForAll.Option.Chat.WithHistory,
		"Display past conversations when resuming a session.",
	)

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

	if err := setBasicChatFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}
	flagSet.BoolVar(
		&aiForAll.Option.Chat.WithHistory,
		"H",
		aiForAll.Option.Chat.WithHistory,
		"Display past conversations when resuming a session.",
	)

	return &ResumeCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetListCommand() (Command, error) {
	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	flagSet.IntVar(
		&aiForAll.Option.List.Count,
		"n",
		aiForAll.Option.List.Count,
		"Print count sessions.",
	)
	flagSet.BoolVar(
		&aiForAll.Option.List.OrderByModify,
		"t",
		aiForAll.Option.List.OrderByModify,
		"Sort by descending time modified (most recently session first).",
	)

	return &ListCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func setBasicChatFlags(aiForAll *AIForAll, flagSet *flag.FlagSet) error {
	flagSet.StringVar(
		&aiForAll.Message,
		"p",
		aiForAll.Message,
		"Message as initial prompt.",
	)
	flagSet.StringVar(
		&aiForAll.Option.Chat.UserPromptTemplate,
		"u",
		aiForAll.Option.Chat.UserPromptTemplate,
		"Name of user prompt template.",
	)
	flagSet.BoolVar(
		&aiForAll.Option.Chat.Interactive,
		"I",
		aiForAll.Option.Chat.Interactive,
		"Runs in interactive mode; set to false when standard input is passed.",
	)
	flagSet.BoolVar(
		&aiForAll.Option.Chat.Stream,
		"S",
		aiForAll.Option.Chat.Stream,
		"Runs in stream mode.",
	)
	flagSet.StringVar(
		&aiForAll.Option.Chat.RunsOn,
		"R",
		aiForAll.Option.Chat.RunsOn,
		"Resume based on the identifier of latest session. (default \"$PPID\")",
	)
	flagSet.BoolVar(
		&aiForAll.Option.Chat.Viewer,
		"V",
		aiForAll.Option.Chat.Viewer,
		fmt.Sprintf(
			"Use the viewer program. (\"%s $SOCK_ADDR\")",
			strings.Join(aiForAll.Option.Chat.ViewerCommand, " "),
		),
	)

	if hasStdin() {
		inputStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		aiForAll.MessageStdin = string(inputStdin)
		aiForAll.Option.Chat.Interactive = false
	}

	return nil
}

func hasStdin() bool {
	if stat, err := os.Stdin.Stat(); err == nil {
		return (stat.Mode() & os.ModeCharDevice) == 0
	}
	return false
}

func workSpaceNotExistError() error {
	return fmt.Errorf("No workspace exists. Please run \"afa init\".")
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
	)
}
