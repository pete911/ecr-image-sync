package docker

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var ecrRegistryRegex *regexp.Regexp

func init() {

	r, err := regexp.Compile(`\d{12}\.dkr\.ecr\.[^.]+\.amazonaws.com`)
	if err != nil {
		log.Fatalf("cannot compile ecr registry regex")
	}
	ecrRegistryRegex = r
}

type Image struct {
	Registry   string
	Repository string
	Tag        string
}

func ToImage(image string) (Image, error) {

	if strings.HasPrefix(image, ":") || strings.HasSuffix(image, ":") {
		return Image{}, fmt.Errorf("image name %s is invalid", image)
	}

	// get tag
	var img Image
	imageTagParts := strings.Split(image, ":")
	if len(imageTagParts) != 1 {
		img.Tag = strings.Join(imageTagParts[1:], ":")
	}

	// get repository if it exists
	imageWithoutVersion := imageTagParts[0]
	imageParts := strings.Split(imageWithoutVersion, "/")
	if len(imageParts) > 1 &&
		(strings.HasPrefix(imageParts[0], "localhost") || strings.Contains(imageParts[0], ".")) {
		// docker.io is default registry and is not included in repository names, if we included this
		// then it is hard to get image id (image id has to be searched by repository name)
		if imageParts[0] != "docker.io" {
			img.Registry = imageParts[0]
		}
		img.Repository = strings.Join(imageParts[1:], "/")
		return img, nil
	}

	img.Repository = imageWithoutVersion
	return img, nil
}

func (i Image) IsECRRegistry() bool {
	return ecrRegistryRegex.MatchString(i.Registry)
}

func (i Image) ECRRegistryAccount() string {

	if i.IsECRRegistry() {
		return strings.Split(i.Registry, ".")[0]
	}
	return ""
}

func (i Image) ECRRegistryRegion() string {

	if i.IsECRRegistry() {
		return strings.Split(i.Registry, ".")[3]
	}
	return ""
}

func (i Image) ToECRImage(awsAccount, awsRegion string) Image {

	awsRegistry := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", awsAccount, awsRegion)
	return Image{
		Registry:   awsRegistry,
		Repository: i.Repository,
		Tag:        i.Tag,
	}
}

func (i Image) String() string {

	image := i.Repository
	if i.Tag != "" {
		image = fmt.Sprintf("%s:%s", image, i.Tag)
	}
	if i.Registry != "" {
		image = fmt.Sprintf("%s/%s", i.Registry, image)
	}
	return image
}
