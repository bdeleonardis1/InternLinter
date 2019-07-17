package main

import (
	"fmt"

	"github.com/bdeleonardis1/InternLinter/config"
	"github.com/bdeleonardis1/InternLinter/github"
	"github.com/bdeleonardis1/InternLinter/linter"
)

func main() {
	args := config.GetArgs()
	configPath, ok := args["--config"]
	if !ok {
		configPath = ""
	}
	config, err := config.GetConfig(configPath)
	if err != nil {
		panic(err)
	}

	branch := github.GetCurrentBranch()
	config.Github.Branch = branch
	title, ok := args["--title"]
	if !ok {
		panic("To open a pull request there must be a title")
	}
	config.Github.Title = title
	diff, err := github.GetDiff(config.Github.Base, config.Github.Branch)
	if err != nil {
		panic(err)
	}
	todos, prints, err := linter.FindFlaws(diff)
	if err != nil {
		panic(err)
	}
	displayProblems(config, todos, prints)
	response, err := github.OpenPrIfNecessary(config, todos, prints)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}

func displayProblems(config *config.Config, todos map[int]string, prints map[int]string) {
	if config.CheckForTODOs && len(todos) > 0 {
		fmt.Println("TODO comments:")
		for line, text := range todos {
			fmt.Println("Line:", line, "-", text)
		}
	}
	fmt.Println()
	if config.CheckForPrints && len(prints) > 0 {
		fmt.Println("Print Statements:")
		for line, text := range prints {
			fmt.Println("Line:", line, "-", text)
		}
	}
	fmt.Println()
}
