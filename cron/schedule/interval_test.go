package schedule

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestByFreq(t *testing.T) {
	type args struct {
		repeat bool
		ctx    []context.Context
	}
	tests := []struct {
		name string
		args args
		want *Interval
	}{
		{
			name: "test by freq",
			args: args{
				repeat: true,
				ctx:    []context.Context{ctx},
			},
			want: &Interval{
				schedule{
					repeat:      true,
					timer:       now(),
					context:     ctx,
					dur:         &duration{location: time.UTC},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByFreq(tt.args.repeat, tt.args.ctx...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByFreq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterval_Next(t *testing.T) {
	type fields struct {
		schedule schedule
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "no timestamp | picks current timestamp",
			fields: fields{
				schedule{
					repeat:      true,
					timer:       now(),
					context:     ctx,
					dur:         &duration{location: time.UTC},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "no date",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Hour:     1,
						Minute:   50,
						Second:   10,
						Nsec:     1000,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: false,
		}, {
			name: "no time",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     1,
						Month:    1,
						Day:      1,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interval{
				schedule: tt.fields.schedule,
			}
			i.AddYear(tt.fields.schedule.dur.Year).
				AddMonth(tt.fields.schedule.dur.Month).
				AddWeek(tt.fields.schedule.dur.Week).
				AddDay(tt.fields.schedule.dur.Day).
				AddHour(tt.fields.schedule.dur.Hour).
				AddMinute(tt.fields.schedule.dur.Minute).
				AddSecond(tt.fields.schedule.dur.Second).
				AddNsec(tt.fields.schedule.dur.Nsec)

			_, err := i.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Interval.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInterval_String(t *testing.T) {
	type fields struct {
		schedule schedule
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no time",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     1,
						Month:    1,
						Day:      1,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			want: "1yrs 1months 0weeks 1days 0hrs 0mins 0secs 0nsecs -> next execution in 0s",
		}, {
			name: "no date",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Hour:     1,
						Minute:   50,
						Second:   10,
						Nsec:     1000,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			want: "0yrs 0months 0weeks 0days 1hrs 50mins 10secs 1000nsecs -> next execution in 0s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interval{
				schedule: tt.fields.schedule,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("Interval.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
