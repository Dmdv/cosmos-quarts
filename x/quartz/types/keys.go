package types

import "errors"

const (
	// ModuleName is the name of the cron module
	ModuleName = "quartz"

	// StoreKey is the default store key for the cron module
	StoreKey = ModuleName

	// RouterKey is the message router key for the cron module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the cron module
	QuerierRoute = ModuleName
)

var (
	// ErrInvalidSchedule is thrown when the input cron schedule is invalid
	ErrInvalidSchedule = errors.New("invalid cron schedule")
)
