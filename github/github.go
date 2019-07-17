package github

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v27/github"
	"github.com/waigani/diffparser"
	"golang.org/x/oauth2"

	"github.com/bdeleonardis1/InternLinter/config"
)

// OpenPrIfNecessary opens a pull request if there are no problems or if the user wants to ignore the problems
func OpenPrIfNecessary(config *config.Config, todos map[int]string, prints map[int]string) (string, error) {
	if (config.CheckForTODOs && len(todos) > 0) || (config.CheckForPrints && len(prints) > 0) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Would you still like to open the PR? (Y/n) ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")
		if text == "Y" {
			return openPR(config)
		}
		return "Okay fix those problems", nil
	}
	return openPR(config)
}

func openPR(config *config.Config) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUBOAUTH")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	head := config.Github.Branch
	if config.Github.IsFork {
		head = config.Github.Username + ":" + head
	}

	newPR := &github.NewPullRequest{
		Title:               github.String(config.Github.Title),
		Head:                github.String(config.Github.Username + ":" + config.Github.Branch),
		Base:                github.String(config.Github.Base),
		MaintainerCanModify: github.Bool(config.Github.MaintainerCanModify),
	}

	pr, _, err := client.PullRequests.Create(context.Background(), config.Github.Organization, config.Github.Repository, newPR)
	if err != nil {
		return "", err
	}

	return pr.GetHTMLURL(), nil
}

// GetCurrentBranch returns what branch you are currently in.
func GetCurrentBranch() string {
	cmd := exec.Command("git", "branch")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "*") {
			branch := strings.Replace(line, "*", "", -1)
			branch = strings.TrimSpace(branch)
			return branch
		}
	}
	return ""
}

// GetDiff returns the diff between the two branches.
func GetDiff(oldBranch string, newBranch string) (*diffparser.Diff, error) {
	cmd := exec.Command("git", "diff", oldBranch+".."+newBranch)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	return diffparser.Parse(out.String())
}
