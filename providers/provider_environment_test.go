package providers

import (
	"context"
	"errors"
	"os"
	"testing"
)

// TestNewDefaultCredsProviderSuccessful tests
func TestNewEnironmentCredsProvider(t *testing.T) {

	os.Setenv("AWS_ACCESS_KEY_ID", "invalid")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "invalid")
	_, err := NewEnvironmentCredsProvider(context.Background(), "def1", WithRegion("us-east-1"))
	if err != nil {
		t.Error("error = nil expected")
	}
}

// TestNewEnvironmentCredsProviderName tests
func TestNewEnvironmentCredsProviderName(t *testing.T) {
	//no AWS Access Key From variable supplied & default is also unset
	_, err := NewEnvironmentCredsProvider(context.Background(), "")
	if err == nil {
		t.Error("expected error != nil, found error == nil")
	} else if !errors.Is(err, ErrInvalidProviderName) {
		t.Errorf("expected: '%v', \nfound: %v", ErrInvalidProviderName.Error(), err.Error())
	}

}

func TestNewEnvironmentCredsProviderAccessKeys(t *testing.T) {

	// no AWS Access Key From variable supplied & default is also unset
	unset("AWS_ACCESS_KEY_ID")
	unset("CUSTOM_ACCESS_KEY_ID")
	_, err := NewEnvironmentCredsProvider(context.Background(), "def1")
	if err == nil {
		t.Error("expected error != nil, found error == nil")
	}

	if err != nil {
		if !errors.Is(err, ErrInvalidAwsAccessKeyIdEnvValue) {
			t.Errorf("expected: '%v', \nfound: %v", ErrInvalidAwsAccessKeyIdEnvValue.Error(), err.Error())
		}
	}

	// Custom AWS Access Key From variable supplied & but all variables unset
	unset("AWS_ACCESS_KEY_ID")
	unset("CUSTOM_ACCESS_KEY_ID")
	_, err = NewEnvironmentCredsProvider(context.Background(), "def1", WithAccessKeyIdFrom("CUSTOM_ACCESS_KEY_ID"))
	if err == nil {
		t.Error("expected error != nil, found error == nil")
	}

	if err != nil {
		if !errors.Is(err, ErrInvalidAwsAccessKeyIdEnvValue) {
			t.Errorf("expected: '%v', \nfound: %v", ErrInvalidAwsAccessKeyIdEnvValue.Error(), err.Error())
		}
	}

	// No AWS Secret Access Key From variable supplied & default is also unset
	os.Setenv("AWS_ACCESS_KEY_ID", "invalid")
	unset("AWS_SECRET_ACCESS_KEY")
	_, err = NewEnvironmentCredsProvider(context.Background(), "def1")

	if err == nil {
		t.Error("expected error != nil, found error == nil")
	}

	if err != nil {

		if !errors.Is(err, ErrInvalidSecretAccessKeyEnvValue) {
			t.Errorf("expected: '%v', \nfound: %v", ErrInvalidSecretAccessKeyEnvValue.Error(), err.Error())
		}
	}

	os.Setenv("AWS_ACCESS_KEY_ID", "invalid")
	unset("AWS_SECRET_ACCESS_KEY")
	_, err = NewEnvironmentCredsProvider(context.Background(), "def1")
	if err == nil {
		t.Error("expected error != nil, found error == nil")
	}

	if err != nil {
		if !errors.Is(err, ErrInvalidSecretAccessKeyEnvValue) {
			t.Errorf("expected: '%v', \nfound: %v", ErrInvalidSecretAccessKeyEnvValue.Error(), err.Error())
		}
	}
}

func TestNewEnvironmentCredsProviderRegion(t *testing.T) {

	os.Setenv("AWS_ACCESS_KEY_ID", "invalid")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "invalid")

	// testing AWS_REGION over defaults
	EUSouth2 := "eu-south-2"
	unset("AWS_REGION")
	unset("CUSTOM_REGION")
	os.Setenv("AWS_REGION", EUSouth2)
	provider, err := NewEnvironmentCredsProvider(context.Background(), "def1")
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != EUSouth2 {
		t.Errorf("expected region '%v', found '%v'", EUSouth2, provider.Config().Region)

	}

	// testing AWS_REGION over WithRegion(r)
	USWest1 := "us-west-1"
	CACentral1 := "ca-central-1"
	unset("AWS_REGION")
	unset("CUSTOM_REGION")
	os.Setenv("AWS_REGION", USWest1)
	provider, err = NewEnvironmentCredsProvider(context.Background(), "def2", WithRegion(CACentral1))
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != USWest1 {
		t.Errorf("expected region '%v', found '%v'", USWest1, provider.Config().Region)

	}

	// testing region from custom env var
	EUCentral1 := "eu-central-1"
	unset("AWS_REGION")
	unset("CUSTOM_REGION")
	os.Setenv("AWS_REGION", USWest1)
	os.Setenv("CUSTOM_REGION", EUCentral1)
	provider, err = NewEnvironmentCredsProvider(context.Background(), "def2", WithRegionFrom("CUSTOM_REGION"), WithRegion(USWest1))
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != EUCentral1 {
		t.Errorf("expected region '%v', found '%v'", EUCentral1, provider.Config().Region)

	}

	// testing region from WithRegion()
	EUSouth1 := "eu-south-1"
	unset("AWS_REGION")
	unset("CUSTOM_REGION")
	provider, err = NewEnvironmentCredsProvider(context.Background(), "def2", WithRegion(EUSouth1))
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != EUSouth1 {
		t.Errorf("expected region '%v', found '%v'", EUSouth1, provider.Config().Region)

	}

	// testing region from WithDefautRegion()
	unset("AWS_REGION")
	unset("CUSTOM_REGION")
	provider, err = NewEnvironmentCredsProvider(context.Background(), "def2", WithDefaultRegion())
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected region '%v', found '%v'", EUSouth1, provider.Config().Region)

	}

	// testing region from defaults
	unset("AWS_REGION")
	unset("CUSTOM_REGION")
	provider, err = NewEnvironmentCredsProvider(context.Background(), "def2")
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected region '%v', found '%v'", EUSouth1, provider.Config().Region)

	}

}

// unset environment variable
func unset(envvar string) {
	os.Setenv(envvar, "")
}
