package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"log"
	"strings"
)

type Client struct {
	Account string
	Region  string
	config  aws.Config
	ecrSvc  *ecr.Client
}

func NewClient(awsAccount, awsRegion string) (Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Client{}, fmt.Errorf("cannot load aws config: %w", err)
	}

	if awsRegion == "" && cfg.Region == "" {
		return Client{}, errors.New("missing aws region")
	}

	if awsRegion != "" {
		cfg.Region = awsRegion
	}
	log.Printf("[info] aws region set to: %s", cfg.Region)

	currentAccount, err := getCurrentAWSAccount(cfg)
	if err != nil {
		return Client{}, fmt.Errorf("get current aws account: %w", err)
	}
	if awsAccount != currentAccount {
		return Client{}, fmt.Errorf("sync account %s is different from current account: %s", awsAccount, currentAccount)
	}

	return Client{Account: awsAccount, Region: cfg.Region, config: cfg, ecrSvc: ecr.NewFromConfig(cfg)}, nil
}

func getCurrentAWSAccount(cfg aws.Config) (string, error) {

	svc := sts.NewFromConfig(cfg)
	resp, err := svc.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return aws.ToString(resp.Account), nil
}

func (c Client) ECRImageExists(repository, imageTag string) (bool, error) {

	req := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repository),
		ImageIds:       []types.ImageIdentifier{{ImageTag: aws.String(imageTag)}},
	}
	resp, err := c.ecrSvc.DescribeImages(context.TODO(), req)
	if err != nil {
		var imageNotFound *types.ImageNotFoundException
		if errors.As(err, &imageNotFound) {
			return false, nil
		}
		var repositoryNotFound *types.RepositoryNotFoundException
		if errors.As(err, &repositoryNotFound) {
			return false, fmt.Errorf("repository %s does not exist, create repository before running ecr-image-sync", repository)
		}
		return false, err
	}

	if len(resp.ImageDetails) != 1 {
		return false, fmt.Errorf("got %d image details, ecpected 1", len(resp.ImageDetails))
	}
	return true, nil
}

func (c Client) GetAuthorizationData() (AuthorizationData, error) {

	resp, err := c.ecrSvc.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return AuthorizationData{}, fmt.Errorf("get authorization token: %v", err)
	}

	if len(resp.AuthorizationData) != 1 {
		return AuthorizationData{}, fmt.Errorf("got %d tokens, ecpected 1", len(resp.AuthorizationData))
	}
	return toLoginData(resp.AuthorizationData[0])
}

func (c Client) GetAuthorizationDataForRegion(region string) (AuthorizationData, error) {

	cfg := c.config.Copy()
	cfg.Region = region

	resp, err := ecr.NewFromConfig(cfg).GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return AuthorizationData{}, fmt.Errorf("get authorization token for %s region: %v", region, err)
	}

	if len(resp.AuthorizationData) != 1 {
		return AuthorizationData{}, fmt.Errorf("got %d tokens, expected 1", len(resp.AuthorizationData))
	}
	return toLoginData(resp.AuthorizationData[0])
}

type AuthorizationData struct {
	User  string
	Token string
	Url   string
}

func toLoginData(data types.AuthorizationData) (AuthorizationData, error) {

	decodedToken, err := base64.StdEncoding.DecodeString(aws.ToString(data.AuthorizationToken))
	if err != nil {
		return AuthorizationData{}, fmt.Errorf("decode token: %v", err)
	}
	token := strings.TrimPrefix(string(decodedToken), "AWS:")

	return AuthorizationData{
		User:  "AWS",
		Token: token,
		Url:   aws.ToString(data.ProxyEndpoint),
	}, nil
}
