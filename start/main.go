package main

import (
	"context"
	"fmt"
	"log"
	"temporalio-example/app"

	"go.temporal.io/sdk/client"
)

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "hugo-netlify-updater-workflow",
		TaskQueue: app.UpdaterTaskQueue,
	}

	// Start the Workflow
	conf, err := app.ParseConfigFile("config.toml")
	if err != nil {
		log.Fatalln("unable to parse config file")
	}

	sourceRepo := conf.SourceRepoReleases
	we, err := c.ExecuteWorkflow(context.Background(), options, app.UpdaterWorkflow, sourceRepo)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var changes int
	err = we.Get(context.Background(), &changes)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}
	printResults(changes, we.GetID(), we.GetRunID())
}

func printResults(changes int, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\nChanges: %d\n\n", changes)
}
