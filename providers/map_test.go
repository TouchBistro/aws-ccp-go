package providers

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {

	name := "def1"

	_, err := NewDefaultCredsProvider(context.Background(), name, WithRegion("us-east-1"))
	if err != nil {
		t.Error("error = nil expected")
	}

	_, err = Get(name)
	if err != nil {
		t.Error("error = nil expected")
	}

}

func TestClone(t *testing.T) {

	name := "def1"
	clone := "def2"

	def1, err := NewDefaultCredsProvider(context.Background(), name, WithRegion("us-east-1"))
	if err != nil {
		t.Error("error = nil expected")
	}

	def2, err := Clone(name, clone)
	if err != nil {
		t.Error("error = nil expected")
	}

	if def1 != def2 {
		t.Error("cloned provider must be the same")
	}

	if MustGet(name) != MustGet(clone) {
		t.Error("MustGet not returning the expected providers")
	}

}
