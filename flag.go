package main

import (
	"flag"
	"os"
	"strconv"
)

type Flags struct {
	ImagesFile string
	AWSAccount string
	AWSRegion  string
	DryRun     bool
	Version    bool
	Args       []string
}

func ParseFlags() (Flags, error) {

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	imagesFile := f.String("images-file", getStringEnv("ECR_SYNC_IMAGES_FILE", "images-list"),
		"file containing list of images to sync")
	awsAccount := f.String("aws-account", getStringEnv("ECR_SYNC_AWS_ACCOUNT", ""),
		"AWS account where to sync images to")
	awsRegion := f.String("aws-region", getStringEnv("ECR_SYNC_AWS_REGION", ""),
		"AWS region where to sync images to")
	dryRun := f.Bool("dry-run", getBoolEnv("ECR_SYNC_DRY_RUN", false),
		"whether to do sync - tag and push image")
	version := f.Bool("version", getBoolEnv("ECR_SYNC_VERSION", false),
		"ecr-image-sync version")

	if err := f.Parse(os.Args[1:]); err != nil {
		return Flags{}, err
	}

	return Flags{
		ImagesFile: stringValue(imagesFile),
		AWSAccount: stringValue(awsAccount),
		AWSRegion:  stringValue(awsRegion),
		DryRun:     boolValue(dryRun),
		Version:    boolValue(version),
		Args:       f.Args(),
	}, nil
}

func getStringEnv(envName string, defaultValue string) string {

	env, ok := os.LookupEnv(envName)
	if !ok {
		return defaultValue
	}
	return env
}

func stringValue(v *string) string {

	if v == nil {
		return ""
	}
	return *v
}

func getBoolEnv(envName string, defaultValue bool) bool {

	env, ok := os.LookupEnv(envName)
	if !ok {
		return defaultValue
	}

	if intValue, err := strconv.ParseBool(env); err == nil {
		return intValue
	}
	return defaultValue
}

func boolValue(v *bool) bool {

	if v == nil {
		return false
	}
	return *v
}
