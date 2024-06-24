package providers

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestNewAssumeroleCredsProvider(t *testing.T) {

	_, err := NewDefaultCredsProvider(context.Background(), "def1", WithRegion("us-east-1"))
	if err != nil {
		t.Error("error = nil expected")
	}
}

func TestNewAssumeRoleCredsProviderName(t *testing.T) {
	//no AWS Access Key From variable supplied & default is also unset
	roleArn := os.Getenv("AWS_CREDS_PROVIDER_TEST_ASSUME_ROLE")
	_, err := NewAssumeRoleCredsProvider(context.Background(), "", WithRoleArn(roleArn))
	if err == nil {
		t.Error("expected error != nil, found error == nil")
	} else if !errors.Is(err, ErrInvalidProviderName) {
		t.Errorf("expected: '%v', \nfound: %v", ErrInvalidProviderName.Error(), err.Error())
	}

}

func TestNewAssumeRoleCredsProviderRegion(t *testing.T) {

	provider, err := NewAssumeRoleCredsProvider(context.Background(), "r1")
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected region '%v', found '%v'", DefaultAWSRegion, provider.Config().Region)

	}

	provider, err = NewAssumeRoleCredsProvider(context.Background(), "r1", WithDefaultRegion())
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != DefaultAWSRegion {
		t.Errorf("expected  region '%v', found '%v'", DefaultAWSRegion, provider.Config().Region)

	}

	USWest2 := "us-west-2"
	provider, err = NewAssumeRoleCredsProvider(context.Background(), "r1", WithRegion(USWest2))
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider.Config().Region != USWest2 {
		t.Errorf("expected region '%v', found '%v'", USWest2, provider.Config().Region)
	}

	EUSouth1 := "eu-south-1"
	base, err := NewDefaultCredsProvider(context.Background(), "d1", WithRegion(EUSouth1))
	if err != nil {
		t.Error("error = nil expected")
	}

	//role to assume
	roleArn := os.Getenv("AWS_CREDS_PROVIDER_TEST_ASSUME_ROLE")
	provider, err = NewAssumeRoleCredsProvider(context.Background(), "r1", WithRegion(USWest2),
		WithBaseCredsProvider(base), WithRoleArn(roleArn), ValidateProvider())
	if err != nil {
		t.Error("error = nil expected")
	}

	if provider == nil {
		t.Error("provider != nil expected")
	}

	if provider.Config().Region != USWest2 {
		t.Errorf("expected region '%v', found '%v'", USWest2, provider.Config().Region)
	}
}
