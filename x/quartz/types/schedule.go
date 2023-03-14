package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"
	"time"
)

// Schedule struct
type Schedule struct {
	Schedule string    `json:"schedule"`
	Height   int64     `json:"height"`
	Msgs     []sdk.Msg `json:"msgs"`
	Labels   []string  `json:"labels"`
}

// NewSchedule creates a new Schedule instance
func NewSchedule(
	schedule string,
	height int64,
	msgs []sdk.Msg,
	labels []string,
) (Schedule, error) {
	if !isValidSchedule(schedule) {
		return Schedule{}, ErrInvalidSchedule
	}

	return Schedule{
		Schedule: schedule,
		Height:   height,
		Msgs:     msgs,
		Labels:   labels,
	}, nil
}

// Next returns the next execution time of the schedule based on the current time
func (s Schedule) Next(currentTime time.Time) (time.Time, error) {
	parser := cronParser{}
	schedule, err := parser.Parse(s.Schedule)
	if err != nil {
		return time.Time{}, ErrInvalidSchedule
	}

	nextTime, _ := schedule.Next(currentTime)
	return nextTime, nil
}

// cronParser struct
type cronParser struct{}

// Parse parses a cron schedule string into a cron schedule
func (p cronParser) Parse(schedule string) (Schedule, error) {
	components := strings.Fields(schedule)
	if len(components) != 5 {
		return Schedule{}, ErrInvalidSchedule
	}

	minute, err := strconv.ParseUint(components[0], 10, 64)
	if err != nil {
		return Schedule{}, ErrInvalidSchedule
	}
	hour, err := strconv.ParseUint(components[1], 10, 64)
	if err != nil {
		return Schedule{}, ErrInvalidSchedule
	}
	dayOfMonth, err := strconv.ParseUint(components[2], 10, 64)
	if err != nil {
		return Schedule{}, ErrInvalidSchedule
	}
	month, err := strconv.ParseUint(components[3], 10, 64)
	if err != nil {
		return Schedule{}, ErrInvalidSchedule
	}
	dayOfWeek, err := strconv.ParseUint(components[4], 10, 64)
	if err != nil {
		return Schedule{}, ErrInvalidSchedule
	}

	scheduleStr := fmt.Sprintf("%d %d %d %d %d", minute, hour, dayOfMonth, month, dayOfWeek)
	return Schedule{Schedule: scheduleStr}, nil
}

// isValidSchedule checks if the input cron schedule is valid
func isValidSchedule(schedule string) bool {
	parser := cronParser{}
	_, err := parser.Parse(schedule)
	if err != nil {
		return false
	}
	return true
}

// ---------------------------------------------------------------------------------------------
