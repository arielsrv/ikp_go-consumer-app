package secrets

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/src/main/app/config/env"
	"github.com/ugurcsen/gods-generic/maps"
	"github.com/ugurcsen/gods-generic/maps/hashmap"
)

type AWSClient interface {
	GetSecretValue(
		ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.GetSecretValueOutput, error)
}

type AWSSecretStore struct {
	AWSClient
	config   aws.Config
	AppCache maps.Map[string, *SecretDto]
}

func NewAWSSecretStore(config aws.Config) *AWSSecretStore {
	return &AWSSecretStore{
		config:    config,
		AWSClient: secretsmanager.NewFromConfig(config),
		AppCache:  hashmap.New[string, *SecretDto](), // non thread safe, need re-deploy
	}
}

func (c AWSSecretStore) Get(key string) *SecretDto {
	secretDto, found := c.AppCache.Get(key)
	if !found {
		secretDto = c.getFromAWS(key)
		if !env.IsEmpty(secretDto.Value) {
			c.AppCache.Put(key, secretDto)
		}
		return secretDto
	}

	return secretDto
}

func (c AWSSecretStore) getFromAWS(key string) *SecretDto {
	secretValueOutput, err := c.GetSecretValue(context.Background(),
		&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(key),
		})

	if err != nil {
		return &SecretDto{
			Err: err,
		}
	}

	secretDto := new(SecretDto)
	secretDto.Key = key
	secretDto.Value = aws.ToString(secretValueOutput.SecretString)

	return secretDto
}
