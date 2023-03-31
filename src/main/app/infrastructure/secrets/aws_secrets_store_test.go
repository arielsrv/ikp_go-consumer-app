package secrets_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/src/main/app/container"
	"github.com/src/main/app/infrastructure/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/ugurcsen/gods-generic/maps/hashmap"
)

type MockSecretStore struct {
	secrets map[string]*string
}

func NewMockSecretStore(values map[string]*string) secrets.AWSSecretStore {
	return secrets.AWSSecretStore{
		AWSClient: &MockSecretStore{
			secrets: values,
		},
		AppCache: hashmap.New[string, *secrets.SecretDto](),
	}
}

func (m *MockSecretStore) GetSecretValue(
	_ context.Context,
	params *secretsmanager.GetSecretValueInput,
	_ ...func(*secretsmanager.Options),
) (*secretsmanager.GetSecretValueOutput, error) {
	if m.secrets[aws.ToString(params.SecretId)] == nil {
		return nil, errors.New("secret not found")
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretString: m.secrets[aws.ToString(params.SecretId)],
	}, nil
}

func TestSecretService_Get(t *testing.T) {
	secretService := NewMockSecretStore(map[string]*string{
		"my_key": aws.String("my_secret"),
	})

	actual := secretService.Get("my_key")
	assert.NoError(t, actual.Err)
	assert.Equal(t, "my_secret", actual.String())

	actual = secretService.Get("my_key")
	assert.NoError(t, actual.Err)
	assert.Equal(t, "my_secret", actual.String())
}

func TestSecretService_GetErr(t *testing.T) {
	secretService := NewMockSecretStore(map[string]*string{})

	actual := secretService.Get("my_key")
	assert.Error(t, actual.Err)
}

func TestNewClient(t *testing.T) {
	secretsService := secrets.NewAWSSecretStore(container.ProvideAWSConfig())
	assert.NotNil(t, secretsService)
}
