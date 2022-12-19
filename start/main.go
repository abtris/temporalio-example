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
	sourceRepo := "gohugoio/hugo"
	we, err := c.ExecuteWorkflow(context.Background(), options, app.UpdaterWorkflow, sourceRepo)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var hugoVersion string
	err = we.Get(context.Background(), &hugoVersion)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}
	printResults(hugoVersion, we.GetID(), we.GetRunID())
}

func printResults(hugoVersion string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", hugoVersion)
}
