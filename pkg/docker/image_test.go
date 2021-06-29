package docker

import (
	"fmt"
	"testing"
)

var (
	testAWSAccount  = "12345678"
	testAWSRegion   = "eu-west-2"
	testAWSRegistry = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", testAWSAccount, testAWSRegion)
)

func TestToImage(t *testing.T) {

	t.Run("image without registry", func(t *testing.T) {
		image, err := ToImage("memcached:1.5")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		if image.Registry != "" {
			t.Errorf("registry: expected empty, actual %s", image.Registry)
		}
		if image.Repository != "memcached" {
			t.Errorf("repository: expected memcached, actual %s", image.Repository)
		}
		if image.Tag != "1.5" {
			t.Errorf("tag: expected 1.5, actual %s", image.Tag)
		}
	})

	t.Run("when image has docker.io registry then registry is removed", func(t *testing.T) {
		image, err := ToImage("docker.io/memcached:1.5")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		if image.Registry != "" {
			t.Errorf("registry: expected empty, actual %s", image.Registry)
		}
		if image.Repository != "memcached" {
			t.Errorf("repository: expected memcached, actual %s", image.Repository)
		}
		if image.Tag != "1.5" {
			t.Errorf("tag: expected 1.5, actual %s", image.Tag)
		}
	})

	t.Run("image with registry", func(t *testing.T) {
		image, err := ToImage("quay.io/jacksontj/promxy:v0.0.58")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		if image.Registry != "quay.io" {
			t.Errorf("registry: expected quay.io, actual %s", image.Registry)
		}
		if image.Repository != "jacksontj/promxy" {
			t.Errorf("repository: expected jacksontj/promxy, actual %s", image.Repository)
		}
		if image.Tag != "v0.0.58" {
			t.Errorf("tag: expected v0.0.58, actual %s", image.Tag)
		}
	})

	t.Run("image ends with : is invalid", func(t *testing.T) {
		_, err := ToImage("quay.io/jacksontj/promxy:")
		if err == nil {
			t.Errorf("no version after ':': expected error %v", err)
		}
	})

	t.Run("image starts with : is invalid", func(t *testing.T) {
		_, err := ToImage(":latest")
		if err == nil {
			t.Errorf("image starts with ':': expected error %v", err)
		}
	})

	t.Run("image name is : is invalid", func(t *testing.T) {
		_, err := ToImage(":")
		if err == nil {
			t.Errorf("image name ':': expected error %v", err)
		}
	})
}

func TestImage_IsECRRegistry(t *testing.T) {

	t.Run("image with ECR registry", func(t *testing.T) {
		image, err := ToImage("602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon/aws-iam-authenticator:v0.5.2-scratch")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		if !image.IsECRRegistry() {
			t.Error("ecr image: expected true, actual false")
		}
	})

	t.Run("image with non-ECR registry", func(t *testing.T) {
		image, err := ToImage("blacklabelops/logrotate:1.3")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		if image.IsECRRegistry() {
			t.Error("ecr image: expected false, actual true")
		}
	})
}

func TestImage_ECRRegistryAccount(t *testing.T) {

	accountNumber := "602401143452"
	region := "us-west-2"

	t.Run("image with ECR registry", func(t *testing.T) {
		image, err := ToImage(fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/amazon/aws-iam-authenticator:v0.5.2-scratch", accountNumber, region))
		if err != nil {
			t.Errorf("ecr image account: unexpected error %v", err)
		}
		if image.ECRRegistryAccount() != accountNumber {
			t.Errorf("ecr image account: expected %s, actual %s", accountNumber, image.ECRRegistryAccount())
		}
	})

	t.Run("image with non-ECR registry", func(t *testing.T) {
		image, err := ToImage("blacklabelops/logrotate:1.3")
		if err != nil {
			t.Errorf("ecr image account: unexpected error %v", err)
		}
		if image.ECRRegistryAccount() != "" {
			t.Errorf("ecr image account: expected empty, actual %s", image.ECRRegistryAccount())
		}
	})
}

func TestImage_ECRRegistryRegion(t *testing.T) {

	accountNumber := "602401143452"
	region := "us-west-2"

	t.Run("image with ECR registry", func(t *testing.T) {
		image, err := ToImage(fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/amazon/aws-iam-authenticator:v0.5.2-scratch", accountNumber, region))
		if err != nil {
			t.Errorf("ecr image region: unexpected error %v", err)
		}
		if image.ECRRegistryRegion() != region {
			t.Errorf("ecr image region: expected %s, actual %s", region, image.ECRRegistryRegion())
		}
	})

	t.Run("image with non-ECR registry", func(t *testing.T) {
		image, err := ToImage("blacklabelops/logrotate:1.3")
		if err != nil {
			t.Errorf("ecr image region: unexpected error %v", err)
		}
		if image.ECRRegistryRegion() != "" {
			t.Errorf("ecr image region: expected empty, actual %s", image.ECRRegistryRegion())
		}
	})
}

func TestToECRImage(t *testing.T) {

	t.Run("image without registry", func(t *testing.T) {
		image, err := ToImage("memcached:1.5")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		ecrImage := image.ToECRImage(testAWSAccount, testAWSRegion)
		if ecrImage.Registry != testAWSRegistry {
			t.Errorf("registry: expected %s, actual %s", testAWSRegistry, ecrImage.Registry)
		}
		if ecrImage.Repository != "memcached" {
			t.Errorf("repository: expected memcached, actual %s", ecrImage.Repository)
		}
		if ecrImage.Tag != "1.5" {
			t.Errorf("tag: expected 1.5, actual %s", ecrImage.Tag)
		}
	})

	t.Run("image without registry and tag", func(t *testing.T) {
		image, err := ToImage("memcached")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		ecrImage := image.ToECRImage(testAWSAccount, testAWSRegion)
		if ecrImage.Registry != testAWSRegistry {
			t.Errorf("registry: expected %s, actual %s", testAWSRegistry, ecrImage.Registry)
		}
		if ecrImage.Repository != "memcached" {
			t.Errorf("repository: expected memcached, actual %s", ecrImage.Repository)
		}
		if ecrImage.Tag != "" {
			t.Errorf("tag: expected empty, actual %s", ecrImage.Tag)
		}
	})

	t.Run("image with registry", func(t *testing.T) {
		image, err := ToImage("quay.io/jacksontj/promxy:v0.0.58")
		if err != nil {
			t.Errorf("registry: unexpected error %v", err)
		}
		ecrImage := image.ToECRImage(testAWSAccount, testAWSRegion)
		if ecrImage.Registry != testAWSRegistry {
			t.Errorf("registry: expected %s, actual %s", testAWSRegistry, ecrImage.Registry)
		}
		if ecrImage.Repository != "jacksontj/promxy" {
			t.Errorf("repository: expected jacksontj/promxy, actual %s", ecrImage.Repository)
		}
		if ecrImage.Tag != "v0.0.58" {
			t.Errorf("tag: expected v0.0.58, actual %s", ecrImage.Tag)
		}
	})
}

func TestToECRImage_String(t *testing.T) {

	t.Run("images without registry are prefixed with aws ecr registry", func(t *testing.T) {

		cases := []struct {
			image    string
			expected string
		}{
			{"memcached:1.5", fmt.Sprintf("%s/memcached:1.5", testAWSRegistry)},
			{"nats-streaming:0.17.0", fmt.Sprintf("%s/nats-streaming:0.17.0", testAWSRegistry)},
			{"hashicorp/terraform:latest", fmt.Sprintf("%s/hashicorp/terraform:latest", testAWSRegistry)},
		}
		for i, c := range cases {
			image, err := ToImage(c.image)
			if err != nil {
				t.Errorf("registry: unexpected error %v", err)
			}
			actual := image.ToECRImage(testAWSAccount, testAWSRegion)
			if actual.String() != c.expected {
				t.Errorf("%d: expected %s, actual %s", i, c.expected, actual)
			}
		}
	})

	t.Run("images with registry are replaced with aws ecr registry", func(t *testing.T) {

		cases := []struct {
			image    string
			expected string
		}{
			{"localhost/test/fluentd:latest", fmt.Sprintf("%s/test/fluentd:latest", testAWSRegistry)},
			{"quay.io/jacksontj/promxy:v0.0.58", fmt.Sprintf("%s/jacksontj/promxy:v0.0.58", testAWSRegistry)},
		}
		for i, c := range cases {
			image, err := ToImage(c.image)
			if err != nil {
				t.Errorf("registry: unexpected error %v", err)
			}
			actual := image.ToECRImage(testAWSAccount, testAWSRegion)
			if actual.String() != c.expected {
				t.Errorf("%d: expected %s, actual %s", i, c.expected, actual)
			}
		}
	})
}
