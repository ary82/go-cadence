package workflows

import (
	"context"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	// timeout for activity task from put in queue to started
	activityScheduleToStartTimeout = time.Second * 10
	// timeout for activity from start to complete
	activityStartToCloseTimeout = time.Minute

	// WorkflowStartToCloseTimeout (from workflow start to workflow close)
	WorkflowStartToCloseTimeout = time.Minute * 20
	// DecisionTaskStartToCloseTimeout (from decision task started to decision task completed, usually very short)
	DecisionTaskStartToCloseTimeout = time.Second * 10
)

// SampleCronResult used to return data from one cron run to next cron run.
type SampleCronResult struct {
	EndTime time.Time
}

// sampleCronWorkflow workflow decider
func CronWorkflow(ctx workflow.Context) (*SampleCronResult, error) {
	workflow.GetLogger(ctx).Info("Cron workflow started.", zap.Time("StartTime", workflow.Now(ctx)))

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: activityScheduleToStartTimeout,
		StartToCloseTimeout:    activityStartToCloseTimeout,
	}
	ctx1 := workflow.WithActivityOptions(ctx, ao)

	startTime := time.Time{} // start from 0 time for first cron job
	if workflow.HasLastCompletionResult(ctx) {
		var lastResult SampleCronResult
		if err := workflow.GetLastCompletionResult(ctx, &lastResult); err == nil {
			startTime = lastResult.EndTime
		}
	}

	endTime := workflow.Now(ctx)

	err := workflow.ExecuteActivity(ctx1, CronActivity, startTime, endTime).Get(ctx, nil)
	if err != nil {
		// cron job failed. but next cron should continue to be scheduled by Cadence server
		workflow.GetLogger(ctx).Error("Cron job failed.", zap.Error(err))
		return nil, err
	}

	return &SampleCronResult{EndTime: endTime}, nil
}

// Cron sample job activity.
func CronActivity(ctx context.Context, beginTime, endTime time.Time) error {
	activity.GetLogger(ctx).Info("Cron job running.", zap.Time("beginTime_exclude", beginTime), zap.Time("endTime_include", endTime))
	return nil
}
