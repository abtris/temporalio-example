package app

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

var client *github.Client

func CheckHugoReleaseVersion(ctx context.Context) (string, error) {
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
	sourceOwner, sourceRepo, err := getRepoPath(conf.SourceRepoReleases)
	if err != nil {
		return "", err
	}
	hugoVersion, releaseURL, releaseInfo, err := getCurrentHugoVersion(ctx, client, sourceOwner, sourceRepo)
	log.Printf("releaseURL: %s\n", releaseURL)
	log.Printf("releaseInfo: %s\n", releaseInfo)
	if err != nil {
		return "", err
	}
	return hugoVersion, nil
}

func CheckCurrentDeployedVersion(ctx context.Context, repository Repository) (DeployResult, error) {
	token := os.Getenv("GITHUB_TOKEN")
	deployResult := DeployResult{}
	if token == "" {
		return deployResult, errors.New("Unauthorized: No token present")
	}
	conf, err := ParseConfigFile("config.toml")
	if err != nil {
		return deployResult, err
	}
	log.Printf("Source repo: %s\n", conf.SourceRepoReleases)
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	owner, repo, err := getRepoPath(repository.Repo)
	if err != nil {
		return deployResult, err
	}
	deployVersion, deployContent, err := getCurrentDeployedVersion(ctx, client, owner, repo, repository.TargetFile, repository.Branch)
	if err != nil {
		return deployResult, err
	}
	log.Printf("deployVersion: %s\n", deployVersion)
	deployResult.deployContent = deployContent
	deployResult.deployVersion = deployVersion
	return deployResult, nil
}

func DeployNewVersion(ctx context.Context, result HugoResult, deployResult DeployResult, repository Repository) (string, error) {
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

	updatedContent := updateVersion(result.hugoVersion, deployResult.deployContent)
	owner, repo, err := getRepoPath(repository.Repo)
	if err != nil {
		return "", err
	}
	commitBranch := getCommitBranch(result.hugoVersion)
	ref, newBranch, err := getRef(ctx, client, owner, repo, repository.Branch, commitBranch)
	if err != nil {
		log.Printf("Error in getRef %v", err)
		return "", err
	}
	if newBranch {
		tree, err := getTree(ctx, client, owner, repo, ref, "netlify.toml", updatedContent)
		if err != nil {
			log.Printf("Error in getTree %v", err)
			return "", err
		}
		errCommit := pushCommit(ctx, client, owner, repo, ref, tree, result.hugoVersion)
		if errCommit != nil {
			log.Printf("Error in pushCommit %v", errCommit)
			return "", errCommit
		}
		errPR := createPullRequest(ctx, client, owner, repo, repository.Branch, result.hugoVersion, result.releaseURL, result.releaseInfo, commitBranch)
		if errPR != nil {
			log.Printf("Error in createPullRequest %v", errPR)
			return "", errPR
		}
	} else {
		log.Printf("PR branch (%s) already exists.\n", commitBranch)
	}
	return "", nil
}
