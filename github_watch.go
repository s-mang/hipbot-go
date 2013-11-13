package main

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"strconv"
	"strings"
	"time"
)

const NO_UPDATES_MESSAGE = "All forks up to date."

type Fork struct {
	Id    int64
	Owner string
	Repo  string
}

func scheduleForkUpdates(delay time.Duration, alertTimeStr string) error {
	alertTime, err := timeToAlert(alertTimeStr)
	if err != nil {
		return err
	}

	// Wait until the starting alert time occurs
	time.Sleep(alertTime.Sub(time.Now()))

	ticker := time.NewTicker(delay)

	go func() {
		for {
			select {
			case <-ticker.C:
				wkday := time.Now().Weekday()
				if wkday != time.Saturday && wkday != time.Sunday {
					speakInHTML(behindForksHTML(), true)
				}
			}
		}
	}()

	return nil
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

	behindHTML := "<strong>A fork update!</strong><br>"

	if len(behindForks) == 0 {
		behindHTML += NO_UPDATES_MESSAGE
	} else {
		behindHTML += "The following forks are out of date:"
		behindHTML += "<ul>"
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

func timeToAlert(alertTimeStr string) (time.Time, error) {
	tim := strings.Split(alertTimeStr, ":")
	if len(tim) != 2 {
		return time.Time{}, errors.New("Time should be of the form hour:min (ie. 14:30)")
	}
	hour, err := strconv.Atoi(tim[0])
	if err != nil {
		return time.Time{}, err
	}
	min, err := strconv.Atoi(tim[1])
	if err != nil {
		return time.Time{}, err
	}
	location := time.Now().Location()
	year, month, day := time.Now().Date()

	alertTime := time.Date(year, month, day, hour, min, 0, 0, location)
	if alertTime.Before(time.Now()) {
		alertTime = alertTime.Add(time.Hour * 24)
	}

	return alertTime, nil
}
