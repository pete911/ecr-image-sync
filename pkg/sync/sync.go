package sync

import (
	"fmt"
	"github.com/pete911/ecr-image-sync/pkg/aws"
	"github.com/pete911/ecr-image-sync/pkg/docker"
	"log"
)

type Client struct {
	awsClient    aws.Client
	dockerClient docker.Client
	registryAuth string
}

func NewClient(awsClient aws.Client, dockerClient docker.Client) (Client, error) {

	registryAuth, err := getRegistryAuth(awsClient)
	if err != nil {
		return Client{}, err
	}

	return Client{awsClient: awsClient, dockerClient: dockerClient, registryAuth: registryAuth}, nil
}

func (c Client) SyncImage(image string, dryRun bool) error {

	publicImage, err := docker.ToImage(image)
	if err != nil {
		return err
	}

	ecrImage := publicImage.ToECRImage(c.awsClient.Account, c.awsClient.Region)
	existsInECR, err := c.awsClient.ECRImageExists(ecrImage.Repository, ecrImage.Tag)
	if err != nil {
		return err
	}
	if existsInECR {
		log.Printf("[info] image %s already exists in ECR", ecrImage)
		return nil
	}

	if err := c.pull(publicImage); err != nil {
		return fmt.Errorf("pull %s image: %v", publicImage, err)
	}

	if dryRun {
		log.Print("[info] dry-run - skipping tag and push image")
		return nil
	}

	if err := c.tagAndPush(publicImage.String(), ecrImage.String()); err != nil {
		return fmt.Errorf("tag and push %s image: %v", publicImage, err)
	}
	return nil
}

func (c Client) pull(image docker.Image) error {

	imageName := image.String()
	isPulled, err := c.dockerClient.IsImagePulled(imageName)
	if err != nil {
		return fmt.Errorf("is image pulled: %v", err)
	}
	if isPulled {
		log.Printf("[info] image %s is already pulled", image)
		return nil
	}

	registryAuth, err := c.getImageECRRegistryAuth(image)
	if err != nil {
		return err
	}

	if err := c.dockerClient.ImagePull(registryAuth, imageName); err != nil {
		return fmt.Errorf("pull %s image: %v", image, err)
	}
	log.Printf("[info] image %s pulled", image)
	return nil
}

// return ECR registry auth of the image to be pulled, or empty string if image is not in ECR registry
// public image in ECR still require registry auth, this can be done with any aws credentials
// this way we can also sync 'public' images in ECR to our private ECR
func (c Client) getImageECRRegistryAuth(image docker.Image) (string, error) {

	if !image.IsECRRegistry() {
		return "", nil
	}

	auth, err := c.awsClient.GetAuthorizationDataForRegion(image.ECRRegistryRegion())
	if err != nil {
		return "", fmt.Errorf("get ecr auth data for %s image", image)
	}

	registryAuth, err := docker.GetRegistryAuth(auth.User, auth.Token, image.Registry)
	if err != nil {
		return "", fmt.Errorf("get docker registry auth: %v", err)
	}
	return registryAuth, nil
}

func (c Client) tagAndPush(image, ecrImage string) error {

	if err := c.dockerClient.ImageTag(image, ecrImage); err != nil {
		return fmt.Errorf("tag image with %s tag: %v", ecrImage, err)
	}
	log.Printf("[info] image %s tagged with %s tag", image, ecrImage)

	if err := c.dockerClient.ImagePush(c.registryAuth, ecrImage); err != nil {
		return fmt.Errorf("push %s image to ECR: %v", ecrImage, err)
	}

	log.Printf("[info] image %s tag pushed to ECR", ecrImage)
	return nil
}

func getRegistryAuth(awsClient aws.Client) (string, error) {

	auth, err := awsClient.GetAuthorizationData()
	if err != nil {
		return "", fmt.Errorf("get aws authorization data: %v", err)
	}
	registryAuth, err := docker.GetRegistryAuth(auth.User, auth.Token, auth.Url)
	if err != nil {
		return "", fmt.Errorf("get docker registry auth: %v", err)
	}
	return registryAuth, nil
}
