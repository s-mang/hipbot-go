package main

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"strings"
	"time"
)

type Fork struct {
	Id    int64
	Owner string
	Repo  string
}

func registerFork(ownerRepo string) string {
	split := strings.Split(ownerRepo, "/")
	if len(split) != 2 {
		return "Invalid fork name. Must be of the form 'original-owner/repo-name'"
	} else if split[0] == forkOwner {
		return "Invalid fork. Please use the original author when registering a fork."
	}

	fork := Fork{}
	DB.Where("owner = ? AND repo = ?", split[0], split[1]).Find(&fork)

	if fork.Id != int64(0) {
		return "Already watching " + ownerRepo
	}

	fork = Fork{Owner: split[0], Repo: split[1]}
	DB.Save(&fork)

	return "Fork successfully registered."
}

func behindForksHTML() string {
	var forks []Fork
	DB.Find(&forks)

	client := github.NewClient(nil)
	behindForks := []string{}

	for _, r := range forks {
		behind, err := r.isBehind(client)
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		} else if behind {
			behindForks = append(behindForks, fmt.Sprintf("<li>%s/%s</li>", r.Owner, r.Repo))
		}
	}

	var behindHTML string

	if len(behindForks) == 0 {
		behindHTML = "All forks up to date."
	} else {
		behindHTML = "<strong>Out-of-date Forks</strong><ul>"
		behindHTML += strings.Join(behindForks, "")
		behindHTML += "</ul>"
	}

	return behindHTML
}

func (f Fork) isBehind(client *github.Client) (bool, error) {
	opt := github.CommitsListOptions{Author: f.Owner}
	upstreamCommits, _, err := client.Repositories.ListCommits(forkOwner, f.Repo, &opt)
	if err != nil {
		return false, err
	} else if len(upstreamCommits) == 0 {
		return false, errors.New("No commits found")
	}
	lastUpstreamCommitTime := *upstreamCommits[0].Commit.Author.Date

	opt = github.CommitsListOptions{Since: lastUpstreamCommitTime.Add(time.Second)}
	newCommits, _, err := client.Repositories.ListCommits(f.Owner, f.Repo, &opt)

	return (len(newCommits) > 0), err
}

func listWatchingForks() string {
	var forks []Fork
	DB.Find(&forks)

	forkedRepos := []string{}
	html := ""

	for _, f := range forks {
		forkedRepos = append(forkedRepos, fmt.Sprintf("<li>%s/%s</li>", f.Owner, f.Repo))
	}

	if len(forkedRepos) == 0 {
		html = "No registered forks."
	} else {
		html = "<strong>Forks Watching</strong><ul>"
		html += strings.Join(forkedRepos, "")
		html += "</ul>"
	}

	return html
}
