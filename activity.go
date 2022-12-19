package app

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

var client *github.Client

func CheckHugoReleaseVersion(ctx context.Context, sourceRepo string) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", errors.New("Unauthorized: No token present")
	}
	conf, err := ParseConfigFile("config.toml")
	if err != nil {
		return "", err
	}
	log.Printf("Source repo: %s\n", conf.SourceRepoReleases)
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	sourceOwner, sourceRepo := getRepoPath(sourceRepo)
	hugoVersion, releaseURL, releaseInfo, err := getCurrentHugoVersion(ctx, client, sourceOwner, sourceRepo)
	log.Printf("releaseURL: %s\n", releaseURL)
	log.Printf("releaseInfo: %s\n", releaseInfo)
	if err != nil {
		return "", err
	}
	return hugoVersion, nil
}

func GetCurrentDeployedVersion(ctx context.Context, repository Repository) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", errors.New("Unauthorized: No token present")
	}
	conf, err := ParseConfigFile("config.toml")
	if err != nil {
		return "", err
	}
	log.Printf("Source repo: %s\n", conf.SourceRepoReleases)
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	owner, repo := getRepoPath(repository.Repo)
	deployVersion, deployContent, err := getCurrentDeployedVersion(ctx, client, owner, repo, repository.TargetFile, repository.Branch)
	if err != nil {
		return "", err
	}
	log.Printf("deployContent: %s\n", deployContent)
	return deployVersion, nil
}
