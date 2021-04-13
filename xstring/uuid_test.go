package xstring

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	t.Log(GenerateID(),randInstance.Int63())
}
