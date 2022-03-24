package schedule

import (
	"context"
	"fmt"
	"time"
)

// Interval is a wrapper for scheduler.
type Interval struct {
	schedule
}

// ByFreq returns a new frequency based scheduler.
func ByFreq(repeat bool, ctx ...context.Context) *Interval {
	return &Interval{
		newSched(repeat, setCtx(ctx)),
	}
}

// AddYear adds years to the scheduler.
//
// eg:
//	...
//	i.AddYear(10)   // will add 10 years to the scheduler.
//	...
func (i Interval) AddYear(val int) Interval {
	return i.add(year, val)
}

// AddMonth adds months to the scheduler.
//
// eg:
//	...
//	i.AddMonth(10)   // will add 10 months to the scheduler.
//	...
func (i Interval) AddMonth(val int) Interval {
	return i.add(month, val)
}

// AddWeek adds weeks to the scheduler.
//
// eg:
//	...
// 	i.AddWeek(2)   // will add 2 weeks to the scheduler.
//	...
func (i Interval) AddWeek(val int) Interval {
	return i.add(week, val)
}

// AddDay adds days to the scheduler.
//
// eg:
//	...
// 	i.AddDay(3)   // will add 3 days to the scheduler.
//	...
func (i Interval) AddDay(val int) Interval {
	return i.add(day, val)
}

// AddHour adds hours to the scheduler.
//
// eg:
//	...
// 	i.AddHour(3)   // will add 3 hours to the scheduler.
//	...
func (i Interval) AddHour(val int) Interval {
	return i.add(hour, val)
}

// AddMinute adds minutes to the scheduler.
//
// eg:
//	...
// 	i.AddMinute(10)   // will add 10 minutes to the scheduler.
//	...
func (i Interval) AddMinute(val int) Interval {
	return i.add(minute, val)
}

// AddSecond adds seconds to the scheduler.
//
// eg:
//	...
//	i.AddSecond(10)   // will add 10 seconds to the scheduler.
//	...
func (i Interval) AddSecond(val int) Interval {
	return i.add(second, val)
}

// AddNsec adds nano seconds to the scheduler.
//
// eg:
//	...
//	i.AddNsec(1000)   // will add 1000 nano seconds to the scheduler.
//	...
func (i Interval) AddNsec(val int) Interval {
	return i.add(nsec, val)
}

// add adds a new interval to the scheduler based on the time unit.
// same time unit can be added multiple times in which case will add up the values to
// determine the next schedule time.
//
// eg:
//	i.add(Year, 10)   // will add 10 years to the scheduler.
// 	i.add(Week, 2)    // will add 2 weeks to the scheduler.
// 	i.add(Nsec, 1000) // will add 1000 nano seconds to the scheduler.
// same can be achieved by chaining the calls for convinience :
//	i.add(Year, 10).add(Week, 2).add(Nsec, 1000) // will add 10 years, 2 weeks and 1000 nano seconds to the scheduler.
//
func (i Interval) add(units timeUnit, val int) Interval {
	i.dur.set(units, val)
	return i
}

// set sets the duration for the scheduler based on the time unit.
//
func (d *duration) set(units timeUnit, val int) {
	switch units {
	case year:
		d.Year += val
	case month:
		d.Month += val
	case week:
		// weeks, 7 days, added to date instead of days to avoid overflow
		d.Week += val
		d.date += val * 7
	case day:
		d.Day += val
	case hour:
		d.Hour += val
	case minute:
		d.Minute += val
	case second:
		d.Second += val
	case nsec:
		d.Nsec += val
	}
}

// Next finds the next scheduler interval the scheduler and prepare for run
func (i Interval) Next() (Scheduler, error) {

	// calculate duration to schedule for
	if dur, err := i.dur.timeUntil(now()); err != nil {
		return nil, err
	} else {
		i.interval = time.Unix(0, dur).In(i.dur.location).Sub(now())
		return &i, nil
	}
}

// timeUntil the duration to schedule for.
func (d *duration) timeUntil(nextSched time.Time) (dur int64, err error) {

	// add time till next schedule in years, months and days
	nextSched = nextSched.AddDate(d.Year, d.Month, d.Day+d.date)

	// add duration till next schedule in hours, minutes, seconds and nanoseconds
	nextSched = nextSched.
		Add(time.Duration(d.Hour) * time.Hour).
		Add(time.Duration(d.Minute) * time.Minute).
		Add(time.Duration(d.Second) * time.Second).
		Add(time.Duration(d.Nsec) * time.Nanosecond).In(d.location)

	// check if the time is in the past
	if !nextSched.After(now()) {
		return dur, fmt.Errorf("time is in the past")
	}
	return nextSched.UnixNano(), err
}

// String returns the string representation of the duration.
// This indicates how often the scheduler will run.
// eg:
//	...
//	i.AddYear(10).AddMonth(10).AddWeek(2).AddDay(3).AddHour(3).AddMinute(10).AddSecond(10).AddNsec(1000)
//	fmt.Println(i.String()) // will print: 10yrs 10months 2weeks 3days 3hrs 10mins 10secs 1000nsecs.
//	...
// ie, it will run every (10 years, 10 months, 2 weeks, 3 days, 3 hours, 10 minutes, 10 seconds and 1000 nano seconds)
// until the scheduler is stopped.
func (i *Interval) String() string {
	return fmt.Sprintf(
		"%dyrs %dmonths %dweeks %ddays %dhrs %dmins %dsecs %dnsecs -> next execution in %s",
		i.dur.Year, i.dur.Month, i.dur.Week, i.dur.Day, i.dur.Hour, i.dur.Minute, i.dur.Second, i.dur.Nsec, i.interval)
}
