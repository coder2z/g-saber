package xlog

import (
	"errors"
	"testing"
)

func TestAuto(t *testing.T) {
	Info("test",FieldErr(errors.New("error")))
	Warn("test",FieldErr(errors.New("error")))
	Error("test",FieldErr(errors.New("error")))
}
