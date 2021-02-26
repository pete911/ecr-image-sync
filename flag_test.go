package main

import (
	"os"
	"testing"
)

func TestDefaultFlags(t *testing.T) {

	setInput(t, nil, nil)

	flags, err := ParseFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flags.ImagesFile != "images-list" {
		t.Errorf("imagesFile: want images-list, got %s", flags.ImagesFile)
	}
	if flags.AWSAccount != "" {
		t.Errorf("AWSAccount: want '', got %s", flags.AWSAccount)
	}
	if flags.AWSRegion != "" {
		t.Errorf("AWSRegion: want empty, got %s", flags.AWSRegion)
	}
}

func TestFlags(t *testing.T) {

	setInput(t, []string{"flag",
		"--images-file", "test-images",
		"--aws-account", "11223344",
		"--aws-region", "eu-west-2",
	}, nil)

	flags, err := ParseFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flags.ImagesFile != "test-images" {
		t.Errorf("imagesFile: want test-images, got %s", flags.ImagesFile)
	}
	if flags.AWSAccount != "11223344" {
		t.Errorf("AWSAccount: want 11223344, got %s", flags.AWSAccount)
	}
	if flags.AWSRegion != "eu-west-2" {
		t.Errorf("AWSRegion: want eu-west-2, got %s", flags.AWSRegion)
	}
}

func TestFlagsFromEnvVar(t *testing.T) {

	setInput(t, []string{"flag"}, map[string]string{
		"ECR_SYNC_AWS_ACCOUNT": "111222333",
		"ECR_SYNC_AWS_REGION":  "eu-west-1",
	})

	flags, err := ParseFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flags.ImagesFile != "images-list" {
		t.Errorf("imagesFile: want images-list, got %s", flags.ImagesFile)
	}
	if flags.AWSAccount != "111222333" {
		t.Errorf("AWSAccount: want 111222333, got %s", flags.AWSAccount)
	}
	if flags.AWSRegion != "eu-west-1" {
		t.Errorf("AWSRegion: want eu-west-1, got %s", flags.AWSRegion)
	}
}

func TestFlagsOverrideEnvVar(t *testing.T) {

	setInput(t, []string{"flag",
		"--aws-region", "eu-west-2",
	}, map[string]string{
		"ECR_SYNC_AWS_REGION": "eu-west-1",
	})

	flags, err := ParseFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flags.AWSRegion != "eu-west-2" {
		t.Errorf("AWSRegion: want eu-west-2, got %s", flags.AWSRegion)
	}
}

// --- helper functions ---

func setInput(t *testing.T, args []string, env map[string]string) {

	osArgs := os.Args
	rollback := func() {
		os.Args = osArgs
		for k := range env {
			os.Unsetenv(k)
		}
	}

	if args == nil {
		args = []string{"test"}
	}

	os.Args = args
	for k, v := range env {
		os.Setenv(k, v)
	}
	t.Cleanup(func() {
		rollback()
	})
}
