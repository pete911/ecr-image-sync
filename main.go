package main

import (
	"bufio"
	"fmt"
	"github.com/pete911/ecr-image-sync/pkg/aws"
	"github.com/pete911/ecr-image-sync/pkg/docker"
	"github.com/pete911/ecr-image-sync/pkg/sync"
	"log"
	"os"
	"strings"
)

var Version = "dev"

func main() {

	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("[fatal] cannot parse flags: %v", err)
	}

	if flags.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	log.Printf("[info] started sync with flags: %+v", flags)

	syncCmd(flags)
}

func syncCmd(flags Flags) {

	awsClient, err := aws.NewClient(flags.AWSAccount, flags.AWSRegion)
	if err != nil {
		log.Fatalf("[fatal] cannot create aws client: %v", err)
	}

	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("[fatal] cannot create docker client: %v", err)
	}

	syncClient, err := sync.NewClient(awsClient, dockerClient)
	if err != nil {
		log.Fatalf("[fatal] cannot create sync client: %v", err)
	}

	var errors bool
	for _, image := range listImages(flags.ImagesFile) {
		if err := syncClient.SyncImage(image, flags.DryRun); err != nil {
			errors = true
			log.Printf("[error] cannot sync image: %v", err)
		}
	}

	// fail if one of the images fail
	if errors {
		os.Exit(1)
	}
}

func listImages(fileName string) []string {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("[fatal] cannot open images list file: %v", err)
	}
	defer file.Close()

	duplicateImages := make(map[string]interface{})
	var images []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		image := strings.TrimSpace(scanner.Text())
		// skip empty lines
		if image == "" {
			continue
		}
		// skip duplicate images
		if _, ok := duplicateImages[image]; ok {
			log.Printf("[warn] duplicate image %s, skipping", image)
			continue
		}
		duplicateImages[image] = nil
		images = append(images, image)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("[fatal] cannot list images file: %v", err)
	}
	return images
}
