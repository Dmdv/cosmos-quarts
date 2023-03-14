package types

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/robfig/cron"
)

// Declarations -----------------------------------------------------------

type Function func(ctx sdk.Context, params map[string]interface{})

// TaskFunc struct
//type TaskFunc struct {
//	Name     string
//	Function func(ctx sdk.Context, params map[string]interface{})
//}

type CronScheduler struct {
	// Internal state and configuration
	cron *cron.Cron
	ctx  sdk.Context
}

var (
	// ErrInvalidSchedule is thrown when the input cron schedule is invalid
	ErrInvalidSchedule = errors.New("invalid cron schedule")
)

// Create a new CronScheduler instance --------------------------------------

// NewCronScheduler creates a new CronScheduler instance
func NewCronScheduler(ctx sdk.Context) CronScheduler {
	return CronScheduler{
		cron: cron.New(),
		ctx:  ctx,
	}
}

// Public API ---------------------------------------------------------------

// RegisterTask Registers a task with the scheduler, using the specified cron syntax and function
func (cs *CronScheduler) RegisterTask(
	spec string,
	taskFn Function,
	settings map[string]interface{},
) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	scheduler, err := parser.Parse(spec)
	if err != nil {
		return ErrInvalidSchedule
	}

	cs.cron.Schedule(scheduler, cron.FuncJob(func() {
		taskFn(cs.ctx, settings)
	}))

	return nil
}

// Lifecycle ---------------------------------------------------------------

// Run The scheduler, executing any tasks that are due
func (cs *CronScheduler) Run() error {
	cs.cron.Start()
	return nil
}

// Stop stops the CronScheduler
func (cs *CronScheduler) Stop() {
	cs.cron.Stop()
}
