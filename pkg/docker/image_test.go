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
		image := ToImage("memcached:1.5")
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
		image := ToImage("docker.io/memcached:1.5")
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
		image := ToImage("quay.io/jacksontj/promxy:v0.0.58")
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
}

func TestImage_IsECRRegistry(t *testing.T) {

	t.Run("image with ECR registry", func(t *testing.T) {
		image := ToImage("602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon/aws-iam-authenticator:v0.5.2-scratch")
		if !image.IsECRRegistry() {
			t.Error("ecr image: expected true, actual false")
		}
	})

	t.Run("image with non-ECR registry", func(t *testing.T) {
		image := ToImage("blacklabelops/logrotate:1.3")
		if image.IsECRRegistry() {
			t.Error("ecr image: expected false, actual true")
		}
	})
}

func TestImage_ECRRegistryAccount(t *testing.T) {

	accountNumber := "602401143452"
	region := "us-west-2"

	t.Run("image with ECR registry", func(t *testing.T) {
		image := ToImage(fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/amazon/aws-iam-authenticator:v0.5.2-scratch", accountNumber, region))
		if image.ECRRegistryAccount() != accountNumber {
			t.Errorf("ecr image account: expected %s, actual %s", accountNumber, image.ECRRegistryAccount())
		}
	})

	t.Run("image with non-ECR registry", func(t *testing.T) {
		image := ToImage("blacklabelops/logrotate:1.3")
		if image.ECRRegistryAccount() != "" {
			t.Errorf("ecr image account: expected empty, actual %s", image.ECRRegistryAccount())
		}
	})
}

func TestImage_ECRRegistryRegion(t *testing.T) {

	accountNumber := "602401143452"
	region := "us-west-2"

	t.Run("image with ECR registry", func(t *testing.T) {
		image := ToImage(fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/amazon/aws-iam-authenticator:v0.5.2-scratch", accountNumber, region))
		if image.ECRRegistryRegion() != region {
			t.Errorf("ecr image region: expected %s, actual %s", region, image.ECRRegistryRegion())
		}
	})

	t.Run("image with non-ECR registry", func(t *testing.T) {
		image := ToImage("blacklabelops/logrotate:1.3")
		if image.ECRRegistryRegion() != "" {
			t.Errorf("ecr image region: expected empty, actual %s", image.ECRRegistryRegion())
		}
	})
}

func TestToECRImage(t *testing.T) {

	t.Run("image without registry", func(t *testing.T) {
		ecrImage := ToImage("memcached:1.5").ToECRImage(testAWSAccount, testAWSRegion)
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
		ecrImage := ToImage("memcached").ToECRImage(testAWSAccount, testAWSRegion)
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
		ecrImage := ToImage("quay.io/jacksontj/promxy:v0.0.58").ToECRImage(testAWSAccount, testAWSRegion)
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
			actual := ToImage(c.image).ToECRImage(testAWSAccount, testAWSRegion)
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
			actual := ToImage(c.image).ToECRImage(testAWSAccount, testAWSRegion)
			if actual.String() != c.expected {
				t.Errorf("%d: expected %s, actual %s", i, c.expected, actual)
			}
		}
	})
}
