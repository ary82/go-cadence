package main

import (
	"context"
	"log"
	"time"

	"github.com/ary82/go-cadence/worker"
	"github.com/google/uuid"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/zap"
)

func main() {
	cadenceClient := worker.BuildCadenceClient()
	logger := worker.BuildLogger()

	domain := "test-domain"
	tasklist := "test-worker"
	workflowID := uuid.New().String()
	requestID := uuid.New().String()
	executionTimeout := int32(60)
	closeTimeout := int32(60)

	workflowType := "main.HelloWorldWorkflow"
	input := []byte(`"World"`)
	cron := "* * * * *"

	req := shared.StartWorkflowExecutionRequest{
		Domain:     &domain,
		WorkflowId: &workflowID,
		WorkflowType: &shared.WorkflowType{
			Name: &workflowType,
		},
		TaskList: &shared.TaskList{
			Name: &tasklist,
		},
		Input:                               input,
		ExecutionStartToCloseTimeoutSeconds: &executionTimeout,
		TaskStartToCloseTimeoutSeconds:      &closeTimeout,
		RequestId:                           &requestID,
		CronSchedule:                        &cron,
	}
	// TODO: Conditional cron

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := cadenceClient.StartWorkflowExecution(ctx, &req)
	if err != nil {
		logger.Error("Failed to create workflow", zap.Error(err))
		panic("Failed to create workflow.")
	}

	logger.Info("successfully started HelloWorld workflow", zap.String("runID", resp.GetRunId()))

	time.Sleep(2 * time.Minute)

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	treq := &shared.TerminateWorkflowExecutionRequest{
		Domain: &domain,
		WorkflowExecution: &shared.WorkflowExecution{
			WorkflowId: &workflowID,
		},
		// FirstExecutionRunID: resp.RunId,
	}
	err = cadenceClient.TerminateWorkflowExecution(ctx, treq)
	if err != nil {
		log.Fatal(err)
	}
}
