package providers

import (
	"context"
	"errors"
	"testing"
)

// TestNewDefaultCredsProviderSuccessful tests
func TestNewDefaultCredsProvider(t *testing.T) {

	_, err := NewDefaultCredsProvider(context.Background(), "def1", WithRegion("us-east-1"))
	if err != nil {
		t.Error("error = nil expected")
	}
}

func TestNewDefaultCredsProviderName(t *testing.T) {
	// no AWS Access Key From variable supplied & default is also unset
	_, err := NewEnvironmentCredsProvider(context.Background(), "")
	if err == nil {
		t.Error("expected error != nil, found error == nil")
	} else if !errors.Is(err, ErrInvalidProviderName) {
		t.Errorf("expected: '%v', \nfound: %v", ErrInvalidProviderName.Error(), err.Error())
	}

}

func TestNewDefaultCredsProviderRegion(t *testing.T) {

	provider, err := NewDefaultCredsProvider(context.Background(), "def1")
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected region '%v', found '%v'", DefaultAWSRegion, provider.Config().Region)

	}

	provider, err = NewDefaultCredsProvider(context.Background(), "def2", WithDefaultRegion())
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected region '%v', found '%v'", DefaultAWSRegion, provider.Config().Region)

	}

	custom_region := "us-west-2"
	provider, err = NewDefaultCredsProvider(context.Background(), "def2", WithRegion(custom_region))
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != custom_region {
		t.Errorf("expected region '%v', found '%v'", custom_region, provider.Config().Region)

	}
}
