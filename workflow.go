package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func UpdaterWorkflow(ctx workflow.Context, sourceRepo string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err := workflow.ExecuteActivity(ctx, CheckHugoReleaseVersion, sourceRepo).Get(ctx, &result)

	return result, err
}
