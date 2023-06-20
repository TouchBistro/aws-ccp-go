package providers

import (
	"context"
	"errors"
	"testing"
)

// TestNewDefaultCredsProviderSuccessful tests
func TestNewSharedCredsProvider(t *testing.T) {

	_, err := NewSharedConfigCredsProvider(context.Background(), "s1")
	if err != nil {
		t.Error("error = nil expected")
	}
}

// TestNewSharedConfigCredsProviderName tests
func TestNewSharedConfigCredsProviderName(t *testing.T) {
	_, err := NewSharedConfigCredsProvider(context.Background(), "")
	if err == nil {
		t.Error("expected error != nil, found error == nil")
	} else if !errors.Is(err, ErrInvalidProviderName) {
		t.Errorf("expected: '%v', \nfound: %v", ErrInvalidProviderName.Error(), err.Error())
	}

}

// TestNewSharedConfigCredsProviderRegion tests
func TestNewSharedConfigCredsProviderRegion(t *testing.T) {

	provider, err := NewSharedConfigCredsProvider(context.Background(), "s1")
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected region '%v', found '%v'", DefaultAWSRegion, provider.Config().Region)

	}

	provider, err = NewSharedConfigCredsProvider(context.Background(), "s1", WithDefaultRegion())
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected region '%v', found '%v'", DefaultAWSRegion, provider.Config().Region)

	}

	USWest2 := "us-west-2"
	provider, err = NewSharedConfigCredsProvider(context.Background(), "s1", WithRegion(USWest2))
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != USWest2 {
		t.Errorf("expected region '%v', found '%v'", USWest2, provider.Config().Region)

	}
}
