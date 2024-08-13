package main

import "fmt"

type Command interface {
	Name() string
	Parse([]string) error
	Run() error
}

type NewCommand struct{}

func (c NewCommand) Name() string               { return "new" }
func (c *NewCommand) Parse(args []string) error { return nil }
func (c *NewCommand) Run() error                { fmt.Println("Run new command."); return nil }

type SourceCommand struct{}

func (c SourceCommand) Name() string               { return "source" }
func (c *SourceCommand) Parse(args []string) error { return nil }
func (c *SourceCommand) Run() error                { fmt.Println("Run source command."); return nil }

type ResumeCommand struct{}

func (c ResumeCommand) Name() string               { return "resume" }
func (c *ResumeCommand) Parse(args []string) error { return nil }
func (c *ResumeCommand) Run() error                { fmt.Println("Run resume command."); return nil }

func GetNewCommand() Command {
	return &NewCommand{}
}

func GetSourceCommand() Command {
	return &SourceCommand{}
}

func GetResumeCommand() Command {
	return &ResumeCommand{}
}
