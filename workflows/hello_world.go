package workflows

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func HelloWorldWorkflow(ctx workflow.Context, name string) (string, error) {
	currentState := "started" // This could be any serializable struct.
	err := workflow.SetQueryHandler(ctx, "current_state", func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = "failed to register query handler"
		return "", err
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("helloworld workflow started")

	numberFailed := 0

	var helloworldResult string
	for i := 0; i < 10; i++ {

		currentState = "executing"
		err := workflow.ExecuteActivity(ctx, HelloWorldActivity, name).Get(ctx, &helloworldResult)
		if err != nil {
			logger.Error("Activity failed.", zap.Error(err))
			numberFailed++
			// return "", err
		}
		currentState = "done, waiting for next"
		time.Sleep(1 * time.Second)
	}

	logger.Info("Workflow completed.", zap.String("Result", helloworldResult))
	currentState = "done"

	fmt.Println("Number of failed activities:", numberFailed)

	return helloworldResult, nil
}

func HelloWorldActivity(ctx context.Context, name string) (string, error) {
	n := rand.Intn(11)
	time.Sleep(time.Duration(n/2) * time.Second)
	if n > 5 {
		return "", fmt.Errorf("example error")
	}
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started")
	return "Hello " + name + "!", nil
}
