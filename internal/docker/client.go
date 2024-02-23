package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/registry"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient() (Client, error) {

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return Client{}, err
	}
	return Client{cli: cli}, nil
}

func (c Client) ImagePull(registryAuth, image string) error {

	reader, err := c.cli.ImagePull(context.Background(), image, types.ImagePullOptions{RegistryAuth: registryAuth})
	if err != nil {
		return err
	}

	printOnOneLine(reader)
	return nil
}

func (c Client) IsImagePulled(image string) (bool, error) {

	id, err := c.getImageID(image)
	if err != nil {
		return false, err
	}
	return id != "", nil
}

func (c Client) ImageTag(image, tag string) error {

	id, err := c.getImageID(image)
	if err != nil {
		return fmt.Errorf("get image id: %v", err)
	}
	if id == "" {
		return errors.New("get image id: image id not found")
	}

	// tagging image with existing tag is OK, no need to check if tag exists
	return c.cli.ImageTag(context.Background(), id, tag)
}

func (c Client) ImagePush(registryAuth, image string) error {

	reader, err := c.cli.ImagePush(context.Background(), image, types.ImagePushOptions{RegistryAuth: registryAuth})
	if err != nil {
		return err
	}

	printOnOneLine(reader)
	return nil
}

func (c Client) getImageID(image string) (string, error) {

	filter := filters.NewArgs()
	filter.Add("reference", image)

	imageSummaries, err := c.cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
	if err != nil {
		return "", fmt.Errorf("cannot list images: %v", err)
	}

	for _, imageSummary := range imageSummaries {
		for _, repoTag := range imageSummary.RepoTags {
			if repoTag == image {
				return imageSummary.ID, nil
			}
		}
	}
	return "", nil
}

func GetRegistryAuth(user, token, url string) (string, error) {

	auth := registry.AuthConfig{
		Username:      user,
		Password:      token,
		ServerAddress: url,
	}
	jsonAuth, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(jsonAuth), nil
}
