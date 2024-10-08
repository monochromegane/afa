package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Command interface {
	Name() string
	Description() string
	Default() bool
	Parse([]string) error
	Run() error
}

type InitCommand struct {
	flagSet  *flag.FlagSet
	aiForAll *AIForAll
}

func (c InitCommand) Name() string { return "init" }

func (c InitCommand) Description() string { return "Initialize configuration and cache directories." }

func (c InitCommand) Default() bool { return false }

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

func (c NewCommand) Description() string { return "Initiates a new session." }

func (c NewCommand) Default() bool { return true }

func (c *NewCommand) Parse(args []string) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	c.aiForAll.Files = c.flagSet.Args()

	if hasStdin() {
		inputStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		c.aiForAll.MessageStdin = string(inputStdin)
		c.aiForAll.Option.Chat.Interactive = false
	}

	if c.aiForAll.Option.Chat.DryRun {
		c.aiForAll.Option.Chat.Interactive = false
	}

	if c.aiForAll.Option.Script.Enabled {
		c.aiForAll.Option.SetScriptOptions()
	} else {
		c.aiForAll.Option.Chat.Quote = false
	}

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

func (c SourceCommand) Description() string { return "Continue from a specified session." }

func (c SourceCommand) Default() bool { return false }

func (c *SourceCommand) Parse(args []string) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	c.aiForAll.Files = c.flagSet.Args()

	if hasStdin() {
		inputStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		c.aiForAll.MessageStdin = string(inputStdin)
		c.aiForAll.Option.Chat.Interactive = false
	}

	c.aiForAll.Option.Chat.Save = true

	if c.aiForAll.Option.Script.Enabled {
		c.aiForAll.Option.SetScriptOptions()
	} else {
		c.aiForAll.Option.Chat.Quote = false
	}

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

func (c ResumeCommand) Description() string { return "Continue from the last session." }

func (c ResumeCommand) Default() bool { return false }

func (c *ResumeCommand) Parse(args []string) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	c.aiForAll.Files = c.flagSet.Args()

	if hasStdin() {
		inputStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		c.aiForAll.MessageStdin = string(inputStdin)
		c.aiForAll.Option.Chat.Interactive = false
	}

	c.aiForAll.Option.Chat.Save = true

	if c.aiForAll.Option.Script.Enabled {
		c.aiForAll.Option.SetScriptOptions()
	} else {
		c.aiForAll.Option.Chat.Quote = false
	}

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

func (c ListCommand) Description() string { return "List sessions." }

func (c ListCommand) Default() bool { return false }

func (c *ListCommand) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *ListCommand) Run() error {
	return c.aiForAll.List()
}

type ShowCommand struct {
	flagSet  *flag.FlagSet
	aiForAll *AIForAll
}

func (c ShowCommand) Name() string { return "show" }

func (c ShowCommand) Description() string { return "Show a specified session." }

func (c ShowCommand) Default() bool { return false }

func (c *ShowCommand) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *ShowCommand) Run() error {
	return c.aiForAll.Show()
}

func GetInitCommand() (Command, error) {
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}
	flagSet.BoolVar(
		&aiForAll.Option.Init.NoInteraction,
		"n",
		aiForAll.Option.Init.NoInteraction,
		"Do not ask interactive question.",
	)

	return &InitCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetNewCommand() (Command, error) {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s new", cmdName), flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	if err := setBasicChatFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}
	if err := setBasicViewerFlags(aiForAll, flagSet); err != nil {
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
	flagSet.BoolVar(
		&aiForAll.Option.Chat.DryRun,
		"dry-run",
		aiForAll.Option.Chat.DryRun,
		"Run in dry-run mode. Outputs only the parsed prompt.",
	)
	flagSet.BoolVar(
		&aiForAll.Option.Chat.Save,
		"L",
		aiForAll.Option.Chat.Save,
		"Save session to the log.",
	)

	return &NewCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetSourceCommand() (Command, error) {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s source", cmdName), flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	if err := setBasicChatFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}
	if err := setBasicViewerFlags(aiForAll, flagSet); err != nil {
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
	flagSet.BoolVar(
		&aiForAll.Option.Chat.MockRun,
		"mock-run",
		aiForAll.Option.Chat.MockRun,
		"Run in mock-run mode. Outputs will be fixed to the last response.",
	)

	return &SourceCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetResumeCommand() (Command, error) {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s resume", cmdName), flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	if err := setBasicChatFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}
	if err := setBasicViewerFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}
	flagSet.BoolVar(
		&aiForAll.Option.Chat.WithHistory,
		"H",
		aiForAll.Option.Chat.WithHistory,
		"Display past conversations when resuming a session.",
	)
	flagSet.BoolVar(
		&aiForAll.Option.Chat.MockRun,
		"mock-run",
		aiForAll.Option.Chat.MockRun,
		"Run in mock-run mode. Outputs will be fixed to the last response.",
	)

	return &ResumeCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func GetListCommand() (Command, error) {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s list", cmdName), flag.ExitOnError)
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

func GetShowCommand() (Command, error) {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s show", cmdName), flag.ExitOnError)
	aiForAll, err := newAIForAll()
	if err != nil {
		return nil, err
	}

	if err := setBasicViewerFlags(aiForAll, flagSet); err != nil {
		return nil, err
	}
	flagSet.StringVar(
		&aiForAll.SessionName,
		"l",
		aiForAll.SessionName,
		"Log name of session.",
	)

	return &ShowCommand{
		flagSet:  flagSet,
		aiForAll: aiForAll,
	}, nil
}

func setBasicChatFlags(aiForAll *AIForAll, flagSet *flag.FlagSet) error {
	flagSet.BoolVar(
		&aiForAll.Option.Script.Enabled,
		"script",
		aiForAll.Option.Script.Enabled,
		"Sets a predefined set of options for script execution simultaneously, setting I, H, S, L and V to false.",
	)
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
		"Runs in interactive mode; set to false when standard input is passed or when in dry-run mode.",
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
		&aiForAll.Option.Chat.Quote,
		"Q",
		aiForAll.Option.Chat.Quote,
		"Wraps the output in double quotes, safely escaped for valid string literals.",
	)

	return nil
}

func setBasicViewerFlags(aiForAll *AIForAll, flagSet *flag.FlagSet) error {
	flagSet.BoolVar(
		&aiForAll.Option.Viewer.Enabled,
		"V",
		aiForAll.Option.Viewer.Enabled,
		fmt.Sprintf(
			"Use the viewer program. (\"%s $SOCK_ADDR\")",
			strings.Join(aiForAll.Option.Viewer.Command, " "),
		),
	)
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
	var configDir string
	if xdgConfigHome, err := getXdgHomeDir("XDG_CONFIG_HOME"); err == nil {
		configDir = xdgConfigHome
	} else {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return nil, err
		}
	}

	var cacheDir string
	if xdgCacheHome, err := getXdgHomeDir("XDG_CACHE_HOME"); err == nil {
		cacheDir = xdgCacheHome
	} else {
		cacheDir, err = os.UserCacheDir()
		if err != nil {
			return nil, err
		}
	}

	return NewAIForAll(
		path.Join(configDir, "afa"),
		path.Join(cacheDir, "afa"),
	)
}

func getXdgHomeDir(env string) (string, error) {
	if xdgHome := os.Getenv(env); xdgHome != "" {
		xdgHome = filepath.Clean(xdgHome)
		if strings.HasPrefix(xdgHome, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return strings.Replace(xdgHome, "~", homeDir, 1), nil
		} else {
			return xdgHome, nil
		}
	}
	return "", errors.New("Not found")
}
