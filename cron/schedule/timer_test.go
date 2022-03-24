package schedule

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestByTimestamp(t *testing.T) {
	type args struct {
		repeat bool
		ctx    []context.Context
	}
	tests := []struct {
		name string
		args args
		want *Timer
	}{
		{
			name: "test by timestamp",
			args: args{
				repeat: true,
				ctx:    []context.Context{ctx},
			},
			want: &Timer{
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
			if got := ByTimestamp(tt.args.repeat, tt.args.ctx...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimer_Next(t *testing.T) {
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
						Year:     2020,
						Month:    1,
						date:     2,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: false,
		}, {
			name: "Every",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     Every,
						Month:    Every,
						Day:      Every,
						date:     Every,
						Hour:     Every,
						Minute:   Every,
						Second:   Every,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "validate year overflow",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     10000,
						Month:    1,
						Day:      1,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "validate month overflow",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     2020,
						Month:    13,
						Day:      1,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "validate date overflow",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     2020,
						Month:    1,
						date:     32,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "validate hour overflow",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     2020,
						Month:    1,
						Day:      1,
						Hour:     24,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "validate minute overflow",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     2020,
						Month:    1,
						Day:      1,
						Hour:     1,
						Minute:   60,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "validate second overflow",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     2020,
						Month:    1,
						Day:      1,
						Hour:     1,
						Minute:   0,
						Second:   60,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		}, {
			name: "validate nano second overflow",
			fields: fields{
				schedule{
					repeat:  true,
					timer:   now(),
					context: ctx,
					dur: &duration{
						Year:     2020,
						Month:    1,
						Day:      1,
						Hour:     1,
						Minute:   0,
						Second:   0,
						Nsec:     1e9,
						location: time.UTC,
					},
					schedEveryN: 0,
					interval:    0,
					tick:        nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Timer{
				schedule: tt.fields.schedule,
			}
			tr.SetYear(tt.fields.schedule.dur.Year).
				SetMonth(Month(tt.fields.schedule.dur.Month)).
				SetDay(Weekday(tt.fields.schedule.dur.Day)).
				SetDate(tt.fields.schedule.dur.date).
				SetHour(tt.fields.schedule.dur.Hour).
				SetMinute(tt.fields.schedule.dur.Minute).
				SetSecond(tt.fields.schedule.dur.Second).
				SetNanosecond(tt.fields.schedule.dur.Nsec).
				SetLocation(tt.fields.schedule.dur.location)
			_, err := tr.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Timer.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTimer_String(t *testing.T) {
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
			want: "0001-01-00 00:00:00.0 UTC Monday",
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
			want: "0000-00-00 01:50:10.1000 UTC Sunday",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Timer{
				schedule: tt.fields.schedule,
			}
			if got := tr.String(); got != tt.want {
				t.Errorf("Timer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
