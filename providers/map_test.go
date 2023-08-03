package providers

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {

	_, err := NewDefaultCredsProvider(context.Background(), "def1", WithRegion("us-east-1"))
	if err != nil {
		t.Error("error = nil expected")
	}

}

func TestClone(t *testing.T) {

	def1, err := NewDefaultCredsProvider(context.Background(), "def1", WithRegion("us-east-1"))
	if err != nil {
		t.Error("error = nil expected")
	}

	def2, err := Clone("def1", "def2")
	if err != nil {
		t.Error("error = nil expected")
	}

	if def1 != def2 {
		t.Error("cloned provider must be the same")
	}

}
