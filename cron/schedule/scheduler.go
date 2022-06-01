package schedule

import (
	"context"
	"time"
)

// timeUnit is custom type for calculating interval by time unit.
// eg: Year, Month, Week, Day, Hour, Minute, Second, Nsec
type timeUnit int

// TimeUnit represents the time units for interval calculation.
const (
	year timeUnit = iota
	month
	week
	day
	date
	hour
	minute
	second
	nsec

	// Every is a special value to represent that the scheduler should run every timeunit value.
	// eg: s.SetMonth(schedule.)
	Every   = -1
	maxYear = 9999
)

// Month type is a replacement for time.Month.
// This is to support custom schedule types.
type Month int

// Month represents the months of year.
const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// Weekday type is a replacement for time.Weekday.
// This is to support custom schedule types.
type Weekday int

// Month represents the months of year.
const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Scheduler interface {
	Next() (Scheduler, error)
	String() string
}

// schedule is the main controller of scheduler pkg.
type schedule struct {
	repeat      bool            // if true, scheduler will run forever as long as requirements are met.
	schedEveryN int             // if > 0, scheduler will run every N interval / timestamp whichever is specified.
	interval    time.Duration   // frequency of interval based scheduler.
	timer       time.Time       // time for next schedule.
	tick        *time.Ticker    // ticker for the schedule which keeps the scheduler logic in wait.
	context     context.Context // context for the scheduler. To control scheduler cancel.
	dur         *duration       // duration for the scheduler is a verbose struct with each time unit in raw format.
}

// duration is a custom type in place for time.Duration
type duration struct {
	Year     int
	Month    int
	Week     int
	Day      int
	date     int
	Hour     int
	Minute   int
	Second   int
	Nsec     int
	location *time.Location
}

// now always returns the current time.
var now = func() time.Time { return time.Now().UTC() }

// newSched sets up a new scheduler with given context.
func newSched(repeat bool, ctx context.Context) schedule {
	return schedule{
		repeat:  repeat,
		timer:   now(),
		context: ctx,
		dur: &duration{
			location: time.UTC,
		},
	}
}

// setCtx sets the context for the scheduler.
// if no ctx is provided from user, default ctx is set
func setCtx(ctx []context.Context) context.Context {
	if len(ctx) > 0 {
		return ctx[0]
	}
	return context.Background()
}
