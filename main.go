package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bdeleonardis1/InternLinter/config"
	"github.com/bdeleonardis1/InternLinter/github"
	"github.com/bdeleonardis1/InternLinter/linter"
)

func main() {
	config, err := config.GetConfig("")
	if err != nil {
		return
	}

	branch := github.GetCurrentBranch()
	config.Github.Branch = branch

	diff, err := github.GetDiff(config.Github.Base, config.Github.Branch)
	todos, prints, err := linter.FindFlaws(diff)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if len(todos) > 0 {
			fmt.Println("TODO comments:")
			for line, text := range todos {
				fmt.Println("Line:", line, "-", text)
			}
		}
		if len(prints) > 0 {
			fmt.Println("Print Statements:")
			for line, text := range prints {
				fmt.Println("Line:", line, "-", text)
			}
		}
	}

	if err != nil || len(todos) > 0 || len(prints) > 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Would you still like to open the PR? (Y/n) ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")
		if text == "Y" {
			github.OpenPR(config)
		} else {
			fmt.Println("Okay, fix those problems")
		}
	} else {
		github.OpenPR(config)
	}
}
