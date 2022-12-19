package main

import (
	"log"
	"temporalio-example/app"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, app.UpdaterTaskQueue, worker.Options{})

	w.RegisterWorkflow(app.UpdaterWorkflow)
	w.RegisterActivity(app.CheckHugoReleaseVersion)
	w.RegisterActivity(app.CheckCurrentDeployedVersion)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
