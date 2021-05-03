package xdata

import (
	"testing"
)

func TestNewConsistent(t *testing.T) {
	type args struct {
		nodeNum int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				nodeNum: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConsistent(tt.args.nodeNum);got!=nil {
				t.Logf("NewConsistent() = %v",got)
			}
		})
	}
}