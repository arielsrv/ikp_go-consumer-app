package container

import (
	"context"
	"fmt"
	"sync"

	"github.com/alicebob/miniredis/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/redis/go-redis/v9"
	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/infrastructure/secrets"
	"github.com/src/main/app/log"
)

var (
	kvsOnce   sync.Once
	kvsClient *redis.Client
)

// ProvideKVSClient
// * Get kvs instance name from AppConfig. Priority order is as follows:
// * 1. If prod and property distributed is true then use cache distributed cache to stop consumer
// * 2. If local use cache shared in-memory to stop consumer.
func ProvideKVSClient() *redis.Client {
	kvsOnce.Do(func() {
		distributed := config.TryBool("consumers.distributed", false)
		if distributed && !env.IsLocal() {
			secretStore := ProvideAWSSecretStore()
			cachePassword := secretStore.Get("cache.password")
			if cachePassword.Err != nil {
				log.Fatal(cachePassword.Err)
			}
			kvsClient = redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d",
					config.String("cache.host"),
					config.TryInt("cache.port", 6379)),
				Password: cachePassword.Value,
				DB:       0,
			})
			log.Warnf("start/stop app by distributed cache")
		} else {
			mr, err := miniredis.Run()
			if err != nil {
				log.Fatal(err)
			}
			kvsClient = redis.NewClient(&redis.Options{
				Addr: mr.Addr(),
			})
			log.Warnf("start/stop app by in-memory cache")
		}
	})

	return kvsClient
}

var (
	awsConfigHandlerOnce sync.Once
	awsConfig            aws.Config
)

func ProvideAWSConfig() aws.Config {
	awsConfigHandlerOnce.Do(func() {
		cfg, err := awsCfg.
			LoadDefaultConfig(context.Background(),
				awsCfg.WithEndpointResolverWithOptions(
					aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						return aws.Endpoint{
							URL:           config.String("aws.url"),
							SigningRegion: config.String("aws.region"),
							PartitionID:   config.String("aws.partition"),
						}, nil
					})))

		if err != nil {
			log.Fatal(err)
		}

		awsConfig = cfg
	})

	return awsConfig
}

var (
	secretsStoreOnce sync.Once
	secretsStore     secrets.SecretStore
)

func ProvideAWSSecretStore() secrets.SecretStore {
	secretsStoreOnce.Do(func() {
		secretsStore = secrets.NewAWSSecretStore(ProvideAWSConfig())
	})

	return secretsStore
}
