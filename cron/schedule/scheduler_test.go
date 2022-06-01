package schedule

import (
	"context"
	"reflect"
	"testing"
	"time"
)

var ctx context.Context

func init() {
	ctx = context.Background()
	now()
	now = func() time.Time { return time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC) }
}

func Test_now(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "now",
			want: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSched(t *testing.T) {
	type args struct {
		repeat bool
		ctx    context.Context
	}

	tests := []struct {
		name string
		args args
		want schedule
	}{
		{
			name: "repeat",
			args: args{
				repeat: true,
				ctx:    ctx,
			},
			want: schedule{
				repeat:      true,
				timer:       now(),
				context:     ctx,
				dur:         &duration{location: time.UTC},
				schedEveryN: 0,
				interval:    0,
				tick:        nil,
			},
		}, {
			name: "once",
			args: args{
				repeat: false,
				ctx:    ctx,
			},
			want: schedule{
				repeat:      false,
				timer:       now(),
				context:     ctx,
				dur:         &duration{location: time.UTC},
				schedEveryN: 0,
				interval:    0,
				tick:        nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSched(tt.args.repeat, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSched() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setCtx(t *testing.T) {
	type args struct {
		ctx []context.Context
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "single ctx",
			args: args{
				ctx: []context.Context{ctx},
			},
			want: ctx,
		}, {
			name: "multiple ctx",
			args: args{
				ctx: []context.Context{ctx, context.Background(), context.Background()},
			},
			want: ctx,
		}, {
			name: "no ctx",
			args: args{
				ctx: []context.Context{},
			},
			want: ctx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setCtx(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}
