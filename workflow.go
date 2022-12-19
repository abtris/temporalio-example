package app

import (
	"log"
	"time"

	"go.temporal.io/sdk/workflow"
)

func UpdaterWorkflow(ctx workflow.Context, sourceRepo string) (int, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err := workflow.ExecuteActivity(ctx, CheckHugoReleaseVersion, sourceRepo).Get(ctx, &result)

	conf, err := ParseConfigFile("config.toml")
	if err != nil {
		log.Fatalf("Missing or wrong config.toml - %v", err)
	}
	log.Printf("Source repo: %s\n", conf.SourceRepoReleases)
	var finalResult int
	finalResult = 0
	for _, repository := range conf.TargetRepository {
		var deployedResult string
		err = workflow.ExecuteActivity(ctx, GetCurrentDeployedVersion, repository).Get(ctx, &deployedResult)
		if deployedResult != result {
			// workflow.ExecuteChildWorkflow() - check result vs deployedResult and execute workflow for PR
			finalResult += 1
		}
	}
	return finalResult, err
}
