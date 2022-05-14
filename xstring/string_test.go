package xstring

import "testing"

func TestExpandEnv(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "a",
			args: args{
				"${USER}",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpandEnv(tt.args.s); got == tt.want {
				t.Errorf("ExpandEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}