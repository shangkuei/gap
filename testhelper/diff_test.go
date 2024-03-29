package helper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type equalA struct {
	unexported string
	Exported   string
}

type equalB struct{}

func TestEqual(t *testing.T) {
	type args struct {
		got  any
		want any
		opts []cmp.Option
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		{
			name: "same struct but different unexported field",
			args: args{
				got: equalA{
					unexported: "unexported",
					Exported:   "exported",
				},
				want: equalA{
					unexported: "",
					Exported:   "exported",
				},
				opts: []cmp.Option{cmp.AllowUnexported(equalA{})},
			},
		},
		{
			name: "same struct",
			args: args{
				got: equalA{
					unexported: "unexported",
					Exported:   "exported",
				},
				want: equalA{
					unexported: "unexported",
					Exported:   "exported",
				},
				opts: []cmp.Option{cmp.AllowUnexported(equalA{})},
			},
			wantOk: true,
		},
		{
			name: "different struct",
			args: args{
				got:  equalA{},
				want: equalB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDiff, gotOk := Equal(tt.args.got, tt.args.want, tt.args.opts...)
			t.Log(gotDiff)
			if got, want := gotOk, tt.wantOk; got != want {
				t.Errorf("Equal() got = %v, want %v", got, want)
			}
		})
	}
}
