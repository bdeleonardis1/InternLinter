package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v27/github"
	"github.com/waigani/diffparser"
	"golang.org/x/oauth2"
)

func main() {
	cmd := exec.Command("git", "diff", "eda68e65..96fc2b3")
	//cmd := exec.Command("ls")
	var out bytes.Buffer
	cmd.Stdout = &out
	os.Chdir("/Users/brian.deleonardis/ast/mongoast")
	cmd.Run()
	diff, _ := diffparser.Parse(out.String())

	todos, prints, err := findFlaws(diff)
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
			openPR()
		} else {
			fmt.Println("Okay, fix those problems")
		}
	} else {
		openPR()
	}
}

func findFlaws(diff *diffparser.Diff) (todos map[int]string, prints map[int]string, err error) {
	todos = make(map[int]string)
	prints = make(map[int]string)
	for _, file := range diff.Files {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, file.NewName, nil, parser.ParseComments)
		if err != nil {
			return nil, nil, err
		}

		parsedComments := make(map[int]string)
		for _, comment := range parsedFile.Comments {
			if strings.Contains(comment.Text(), "TODO") {
				parsedComments[fset.Position(comment.Pos()).Line] = strings.TrimSpace(comment.Text())
			}
		}

		parsedPrintlns := make(map[int]string)
		ast.Inspect(parsedFile, func(x ast.Node) bool {
			c, ok := x.(*ast.CallExpr)
			if !ok {
				return true
			}
			s, ok := c.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			i, ok := s.X.(*ast.Ident)
			if !ok {
				return true
			}
			if i.Name == "fmt" && s.Sel.String() == "Println" {
				parsedPrintlns[fset.Position(x.Pos()).Line] = "fmt.Println"
			}
			return false
		})

		for _, hunk := range file.Hunks {
			for _, line := range hunk.NewRange.Lines {
				if line.Mode == diffparser.ADDED {
					lineContent := strings.TrimSpace(line.Content)

					commentText, ok := parsedComments[line.Number]
					if ok && len(lineContent) > 4 && lineContent[0:2] == "//" && commentText == lineContent[3:] {
						todos[line.Number] = lineContent
					}

					_, ok = parsedPrintlns[line.Number]
					if ok && len(lineContent) > 11 && lineContent[0:11] == "fmt.Println" {
						prints[line.Number] = lineContent
					}
				}
			}
		}
	}

	return todos, prints, nil
}

func openPR() (string, error) {
	ctx := context.Background()
	fmt.Println("os.Getenv:", os.Getenv("GITHUBOAUTH"))
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUBOAUTH")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	newPR := &github.NewPullRequest{
		Title:               github.String("Testing PR"),
		Head:                github.String("bdeleonardis1:feature"),
		Base:                github.String("master"),
		Body:                github.String("Testing testing 123"),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(context.Background(), "codebase-berkeley-mentored-project-fa17", "LinterTester", newPR)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return "", nil
}

// If I decide to refactor again I can always use this:

// ProblemType is a type of problem
type ProblemType string

// constants for problem types
const (
	Todo  ProblemType = ProblemType("Todo")
	Print ProblemType = ProblemType("Print")
)

// Problem stores all of the information you need for a problem.
type Problem struct {
	Type ProblemType
	Line int
	Text string
}

// NewProblem creates a problem.
func NewProblem(problemType ProblemType, line int, text string) *Problem {
	return &Problem{
		Type: problemType,
		Line: line,
		Text: text,
	}
}
