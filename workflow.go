package app

import (
	"log"
	"time"

	"go.temporal.io/sdk/workflow"
)

type HugoResult struct {
	hugoVersion string
	releaseURL  string
	releaseInfo string
}

type DeployResult struct {
	deployVersion string
	deployContent string
}

func UpdaterWorkflow(ctx workflow.Context, sourceRepo string) (int, error) {
	var finalResult int
	finalResult = 0
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var result HugoResult
	err := workflow.ExecuteActivity(ctx, CheckHugoReleaseVersion, sourceRepo).Get(ctx, &result)

	conf, err := ParseConfigFile("config.toml")
	if err != nil {
		log.Fatalf("Missing or wrong config.toml - %v", err)
	}
	log.Printf("Source repo: %s\n", conf.SourceRepoReleases)
	for _, repository := range conf.TargetRepository {
		var deployedResult DeployResult
		err = workflow.ExecuteActivity(ctx, CheckCurrentDeployedVersion, repository).Get(ctx, &deployedResult)
		if deployedResult.deployVersion != result.hugoVersion {
			var resultChild bool
			err = workflow.ExecuteActivity(ctx, DeployNewVersion, result, deployedResult, repository).Get(ctx, &resultChild)
			if err != nil {
				return finalResult, err
			}
			finalResult += 1
		}
	}
	return finalResult, err
}
