package schedule

import "context"

// timer is a wrapper for scheduler.
type timer struct {
	schedule
}

// ByTimestamp returns a new clock time based scheduler.
func ByTimestamp(repeat bool, ctx ...context.Context) *timer {
	return &timer{
		newSched(repeat, setCtx(ctx)),
	}
}
