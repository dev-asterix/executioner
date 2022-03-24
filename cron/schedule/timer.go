package schedule

import (
	"context"
	"fmt"
	"time"
)

// Timer is a wrapper for scheduler.
type Timer struct {
	schedule
}

// ByTimestamp returns a new clock time based scheduler.
func ByTimestamp(repeat bool, ctx ...context.Context) *Timer {
	return &Timer{
		newSched(repeat, setCtx(ctx)),
	}
}

// SetYear sets the year of the scheduler.
//
// eg:
//	...
//	t.SetYear(2050)   // will set the year of the scheduler to 2050.
//	...
func (t Timer) SetYear(val int) Timer {
	t.set(year, val)
	return t
}

// SetMonth sets the month of the scheduler.
//
// eg:
//	...
//	t.SetMonth(time.January)   // will set the month of the scheduler to January of said year.
//	...
func (t Timer) SetMonth(val Month) Timer {
	t.set(month, int(val))
	return t
}

// SetDate sets the date of the scheduler.
//
// eg:
//	...
//	t.SetDate(29)   // will set the date of the scheduler to 29th of the month in given year.
//	...
func (t Timer) SetDate(val int) Timer {
	t.set(date, val)
	return t
}

// SetDay sets the day of the scheduler.
//
// eg:
//	...
//	t.SetDay(time.Friday)   // will set the day of the scheduler to Friday of given month and year.
//	...
func (t Timer) SetDay(val Weekday) Timer {
	t.set(day, int(val))
	return t
}

// SetHour sets the hour of the scheduler.
//
// eg:
//	...
//	t.SetHour(3)   // will set the hour of the scheduler to 3 of the day of month and year.
//	...
func (t Timer) SetHour(val int) Timer {
	t.set(hour, val)
	return t
}

// SetMinute sets the minute of the scheduler.
//
// eg:
//	...
//	t.SetMinute(10)   // will set the minute of the scheduler to 10th min of the hour of day and year.
//	...
func (t Timer) SetMinute(val int) Timer {
	t.set(minute, val)
	return t
}

// SetSecond sets the second of the scheduler.
//
// eg:
//	...
//	t.SetSecond(10)   // will set the second of the scheduler to 10th sec of the minute of hour of day and year.
//	...
func (t Timer) SetSecond(val int) Timer {
	t.set(second, val)
	return t
}

// SetNanosecond sets the nanosecond of the scheduler.
//
// eg:
//	...
//	t.SetNanosecond(10)   // will set the nanosecond of the scheduler to 10th of the second of minute of hour of day and year.
//	...
func (t Timer) SetNanosecond(val int) Timer {
	t.set(nsec, val)
	return t
}

// SetLocation sets the location of the scheduler.
//
// eg:
//	...
//	t.SetLocation(time.UTC)   // will set the location of the scheduler to UTC.
//	...
func (t Timer) SetLocation(loc *time.Location) Timer {
	t.dur.location = loc
	return t
}

// set sets the time of execution for the scheduler based on the time unit.
//
func (t *Timer) set(units timeUnit, val int) {
	switch units {
	case year:
		t.dur.Year = val
	case month:
		t.dur.Month = val
	case day:
		t.dur.Day = val
	case date:
		t.dur.date = val
	case hour:
		t.dur.Hour = val
	case minute:
		t.dur.Minute = val
	case second:
		t.dur.Second = val
	case nsec:
		t.dur.Nsec = val
	}
}

// Init the scheduler and prepare for run
func (t Timer) Next() (Scheduler, error) {
	next, err := t.dur.nextDate()
	if err != nil {
		return nil, err
	}
	t.timer = next
	t.tick = time.NewTicker(next.Sub(now()))
	return &t, err
}

// nextDate sets the next date of execution for the scheduler based on the time unit.
func (d *duration) nextDate() (next time.Time, err error) {

	// validate duration for time based scheduler
	if err = d.validate(); err != nil {
		return
	}

	// set the next date of scheduler
	if next, err = d.findNextDate(); err != nil {
		return
	}

	// validate date
	if !next.After(now()) {
		return next, fmt.Errorf("date must be after %s, given %s", now().Format(time.RFC3339), next.Format(time.RFC3339))
	}
	return
}

// findNextDate returns the next date of execution for the scheduler based on the time unit.
func (d *duration) findNextDate() (updatedNext time.Time, err error) {

	next := time.Date(d.Year, time.Month(d.Month), d.date, d.Hour, d.Minute, d.Second, d.Nsec, d.location)
	// durNext is the duration of next date of scheduler
	durNext := duration{
		Year:     next.Year(),
		Month:    int(next.Month()),
		Week:     d.Week,
		Day:      int(next.Weekday()),
		date:     next.Day(),
		Hour:     next.Hour(),
		Minute:   next.Minute(),
		Second:   next.Second(),
		Nsec:     next.Nanosecond(),
		location: next.Location(),
	}

	if d.Year == Every || d.Year < 1 {
		durNext.Year = now().Year()

		// update the scheduler timer, if error is returned, then return the error.
		next, err = durNext.update(year, next)
		if err != nil {
			return
		}
	}
	if d.Month == Every || d.Month < int(January) {
		durNext.Month = int(now().Month())
		durNext.update(month, next)
	}
	if d.Day == Every || d.Day < int(Sunday) {
		durNext.Day = int(now().Weekday())
		durNext.update(day, next)
	}
	if d.date == Every || d.date < 1 {
		durNext.date = now().Day()
		durNext.update(date, next)
	}
	if d.Hour == Every || d.Hour < 0 {
		durNext.Hour = now().Hour()
		durNext.update(hour, next)
	}
	if d.Minute == Every || d.Minute < 0 {
		durNext.Minute = now().Minute()
		durNext.update(minute, next)
	}
	if d.Second == Every || d.Second < 0 {
		durNext.Second = now().Second()
		durNext.update(second, next)
	}
	updatedNext = time.Date(durNext.Year, time.Month(durNext.Month), durNext.date, durNext.Hour, durNext.Minute, durNext.Second, durNext.Nsec, durNext.location)
	return
}

// update updates the next date of scheduler based on the time unit.
// It is used to update the next date of scheduler when the time unit is set to Every.
// If updated time is in the past, it will update the next date of scheduler to next time unit.
func (d *duration) update(units timeUnit, next time.Time) (time.Time, error) {

	var attemptsRem = 100

reschedule: // reschedule reruns the logic until a valid time is found
	attemptsRem--

	// if no more attempts are left, then return the error
	if attemptsRem <= 0 {
		return time.Time{}, fmt.Errorf("unable to find a valid upcoming date which can be scheduled matching given conditions.")
	}

	// set the next date of scheduler
	next = time.Date(d.Year, time.Month(d.Month), d.date, d.Hour, d.Minute, d.Second, d.Nsec, d.location)

	// check if the next date is in the past
	if now().After(next) {
		switch units {
		case year:
			// update the year to next
			d.Year++

			// check overflow
			d.validateYear()
			goto reschedule

		case month:
			// update the month to next
			d.Month++

			// check overflow
			d.validateMonth()
			goto reschedule

		case day:
			// update the day to next
			d.Day++
			goto reschedule

		case date:
			// update the date to next
			d.date++

			// check overflow
			d.validateDate()
			goto reschedule

		case hour:
			// update the hour to next
			d.Hour++

			// check overflow
			d.validateHour()
			goto reschedule

		case minute:
			// update the minute to next
			d.Minute++

			// check overflow
			d.validateMinute()
			goto reschedule

		case second:
			// update the second to next
			d.Second++

			// check overflow
			d.validateSecond()
			goto reschedule

		}
	}
	return next, nil
}

// daysOfMonth returns the number of days in the given month of the year.
func daysOfMonth(month, year int, loc *time.Location) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, loc).AddDate(0, 0, -1).Day()
}

// validateYear validates the year of the given duration. If overflow is found, it will update the year to next.
func (d *duration) validateYear() {
	if d.Year > 9999 {
		d.Year = now().Year()
	}
}

// validateMonth validates the month of the given duration. If overflow is found, it will update the year to next.
func (d *duration) validateMonth() {
	if d.Month > 12 {
		d.Month = 1
		d.Year++
		d.validateYear()
	}
}

// validateDate validates the date of the given duration. If overflow is found, it will update the month to next.
func (d *duration) validateDate() {
	if d.date > daysOfMonth(d.Month, d.Year, d.location) {
		d.date = 1
		d.Month++
		d.validateMonth()
	}
}

// validateHour validates the hour of the given duration. If overflow is found, it will update the date to next.
func (d *duration) validateHour() {
	if d.Hour > 23 {
		d.Hour = 0
		d.date++
		d.validateDate()
	}
}

// validateMinute validates the minute of the given duration. If overflow is found, it will update the hour to next.
func (d *duration) validateMinute() {
	if d.Minute > 59 {
		d.Minute = 0
		d.Hour++
		d.validateHour()
	}
}

// validateSecond validates the second of the given duration. If overflow is found, it will update the minute to next.
func (d *duration) validateSecond() {
	if d.Second > 59 {
		d.Second = 0
		d.Minute++
		d.validateMinute()
	}
}

// validate validates the duration for time based scheduler.
func (d *duration) validate() error {

	// validate year
	if d.Year > maxYear {
		return fmt.Errorf("year must be between %d and %d, where -1 represents every year", Every, maxYear)
	}
	// validate month
	if d.Month > 12 {
		return fmt.Errorf("month must be between 1 and 12 or schedule.Every which represents every month")
	}
	// validate day
	if d.date > 31 {
		return fmt.Errorf("day must be between 1 and 31 or schedule.Every which represents every day")
	}
	// validate hour
	if d.Hour > 23 {
		return fmt.Errorf("hour must be between 0 and 23 or schedule.Every which represents every hour")
	}
	// validate minute
	if d.Minute > 59 {
		return fmt.Errorf("minute must be between 0 and 59 or schedule.Every which represents every minute")
	}
	// validate second
	if d.Second > 59 {
		return fmt.Errorf("second must be between 0 and 59 or schedule.Every which represents every second")
	}
	// validate nanosecond
	if d.Nsec > 999999999 {
		return fmt.Errorf("nanosecond must be between 0 and 999999999 or schedule.Every which represents every nanosecond")
	}
	return nil
}

// String returns the string representation of the duration.
// This indicates how often the scheduler will run.
// eg:
//	...
//   t := schedule.ByTimestamp(true, ctx)
//   t.SetYear(2020).SetMonth(time.January).SetDate(1).SetDay(time.Monday).SetHour(0).SetMinute(0).SetSecond(0).SetNanosecond(0)
//   s := t.String()
//   fmt.Println(s)
//	...
//	// output:
//	// 2020-01-01 00:00:00.000000000 Mon UTC
// ie, it will run only if all the conditions match for the day.
func (t *Timer) String() string {
	return fmt.Sprintf(
		"%04d-%02d-%02d %02d:%02d:%02d.%d %s %s",
		t.dur.Year, t.dur.Month, t.dur.date, t.dur.Hour, t.dur.Minute, t.dur.Second, t.dur.Nsec, t.dur.location, time.Weekday(t.dur.Day),
	)
}
