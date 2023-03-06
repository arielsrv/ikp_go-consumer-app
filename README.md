[![pipeline status](https://gitlab.tiendanimal.com:8088/iskaypet/digital/tools/dev/go-consumer-app/badges/main/pipeline.svg)](https://gitlab.tiendanimal.com:8088/iskaypet/digital/tools/dev/go-consumer-app/-/commits/main)
[![coverage report](https://gitlab.tiendanimal.com:8088/iskaypet/digital/tools/dev/go-consumer-app/badges/main/coverage.svg)](https://gitlab.tiendanimal.com:8088/iskaypet/digital/tools/dev/go-consumer-app/-/commits/main)
[![release](https://gitlab.tiendanimal.com:8088/iskaypet/digital/tools/dev/go-consumer-app/-/badges/release.svg)](https://gitlab.tiendanimal.com:8088/iskaypet/digital/tools/dev/go-consumer-app/-/releases)

> This SDK provides a SaaS framework to build middle-end for consumers and APIs

The intent of the project is to provide a lightweight microservice sdk, based on Golang

The main goal is to provide a modular framework with high level abstractions to receive messages in
any technology, which enforces best
practices

You don't need to know how to handle messages in Amazon SQS. You only need to know HTTP protocol.

```
Send message ──> Topic (Amazon SNS)
                └──> Consumer (Amazon SNS)
                    └──> Pusher (HTTP Client)
                        └──> Your API (HTTP endpoint to receive messages)
```

## Table of contents

* [Layers](#layers)
* [Project setup](#project-setup)
* [SDK](#sdk)
    * [Configuration](#configuration)
        * [AWS](#AWS)
        * [Queues](#Queues)
        * [Consumer](#consumer)
        * [Pusher](#Pusher)
        * [RestClient](#restclient)
            * [RestClient configuration](#restclient-configuration)
            * [RestClient usage](#restclient-usage)
    * [Metrics](#metrics)
        * [Pusher dashboard](#pusher-dashboard)
    * [Contributors](#contributors)
    * [Support](#support)

## Layers

<img width="1211" alt="image" src="https://user-images.githubusercontent.com/760657/220067058-3f86739d-eaa5-49fd-9000-663d99f81642.png">

> Example: 8 messages could be handled by 2 workers and 4 receivers

<img width="499" alt="image" src="https://user-images.githubusercontent.com/760657/221838461-5dc11b39-4c13-43d4-8442-c4174bbc1d24.png">

## Project setup

Local development

```shell
brew install localstack
```

```shell
pip install awscli-local
```

```shell
awslocal sqs create-queue --queue-name orders-consumer
awslocal sqs purge-queue --queue-url http://localhost:4566/000000000000/orders-consumer
awslocal sqs list-queues
awslocal sqs receive-message --queue-url http://localhost:4566/000000000000/orders-consumer
awslocal sns create-topic --name orders-topic
awslocal sns list-subscriptions
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:orders-topic --protocol sqs --notification-endpoint arn:aws:sns:us-east-1:000000000000:orders-consumer
awslocal sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/orders-consumer --attribute-names All
awslocal sns publish --topic-arn arn:aws:sns:us-east-1:000000000000:orders-topic --message '{"order_id": 1}'

```

```shell
brew install go-task/tap/go-task
```

```shell
task build run

```

TODO

<img width="1422" alt="image" src="https://user-images.githubusercontent.com/760657/221419837-025af527-227c-4941-90be-8036a531ad44.png">

## SDK

### Configuration

Environment configuration is based on **Archaius Config**, you should use a similar folder
structure.
*SCOPE* env variable in remote environment is required from Kubernetes

```
└── config
    ├── config.yml (shared config)
    └── local
        └── config.yml (for local development)
    └── prod (for remote environment)
        └── config.yml (base config)
        └── {environment}.config.yml (base config)
```

The SDK provides a simple configuration hierarchy

* env variables
* resources/config/config.properties (shared config)
* resources/config/{environment}/config.properties (override shared config by environment)
* resources/config/{environment}/{scope}.config.properties (override env and shared config by scope)

example *consumers-api.uat.dp.iskaypet.com*

```
└── env variables                               (always first)
└── config
    ├── config.yml                              3th (third)
    └── local
        └── config.yml                          ignored
    └── prod
        └── config.yml (base config)            2nd (second)
        └── uat.config.yml (base config)        1st (first)
```

* 1st (first)   prod/uat.config.yml
* 2nd (second)  prod/config.yml
* 3th (third)   config.yml

```
2023-02-26 17:10:35 [INFO] working directory: /app
2023-02-26 17:10:35 [INFO] loaded configuration file: /app/src/resources/config/prod/uat.config.yml
2023-02-26 17:10:35 [INFO] loaded configuration file: /app/src/resources/config/prod/config.yml
2023-02-26 17:10:35 [INFO] loaded configuration file: /app/src/resources/config/config.yml
2023-02-26 17:10:35 [INFO] invoke dynamic handler:FileSource
2023-02-26 17:10:35 [INFO] enable env source
2023-02-26 17:10:35 [INFO] invoke dynamic handler:EnvironmentSource
2023-02-26 17:10:35 [INFO] archaius init success
2023-02-26 17:10:35 [WARN] ENV: prod, SCOPE: uat
2023-02-26 17:10:35 [WARN] warn: config SCOPE not found, fallback to empty string
2023-02-26 17:10:35 [INFO] create new watcher
2023-02-26 17:10:35 [INFO] Listening on local address 0.0.0.0:8080
2023-02-26 17:10:35 [INFO] Open https://consumers-api.uat.dp.iskaypet.com/ping in the browser
```

#### AWS

```yaml
# aws
aws:
  url: http://localhost:4566
  id: development
  secret: development
  profile: default
  region: us-east-1
```

#### Queues

```yaml
# queues-clients
queues:
  users:
    name: users-consumer
    parallel: 2 # default is 2
    timeout: 1000
```

#### Consumer

Queue to consume messages.

```yaml
# consumers
consumers:
  users:
    workers: 10 # default is instances core - 1
```

#### Pusher

Your app to receive messages. Example: my.app/news. Must allow POST Http Request in
Application/json.

```yaml
# pusher (your-app)
pusher:
  target-endpoint: my.app/news
```

```
POST my.app/news
Content-Type: application/json
```

```json
{
  "id": "message_unique_identifier",
  "msg": "the_json_embedded_message"
}
```

```yaml
# target-app (your-app)
target-app:
  endpoint: my.app/news
```

#### RestClient

Pusher app need a rest client to send messages to target.

##### RestClient configuration

```yaml
# rest-pools
rest:
  pool: # isolated custom-pool for any or specific client
    default:
      pool:
        size: 20
        timeout: 2000
        connection-timeout: 5000
  client: # specific client with default pool
    target-app:
      pool: default
```

##### RestClient usage

```gotemplate
func (c HTTPPusherClient) PostMessage(requestBody *RequestBody) error {
	startTime := time.Now()
	response := c.rb.Post(c.targetEndpoint, requestBody)
	elapsedTime := time.Since(startTime)

	metrics.Collector.RecordExecutionTime(metrics.PusherHTTPTime, elapsedTime.Milliseconds())

	if response.Err != nil {
		var err net.Error
		if ok := errors.As(response.Err, &err); ok && err.Timeout() {
			log.Warnf("pusher timeout, discuss cap theorem, possible inconsistency ensure handle duplicates from target app, MessageId: %s", requestBody.ID)
			metrics.Collector.IncrementCounter(metrics.PusherHTTPTimeout)
		}
		return response.Err
	}

	switch {
	case response.StatusCode >= 200 && response.StatusCode < 300:
		metrics.Collector.IncrementCounter(metrics.PusherStatusOK)
	case response.StatusCode >= 400 && response.StatusCode < 500:
		metrics.Collector.IncrementCounter(metrics.PusherStatus40x)
	case response.StatusCode >= http.StatusInternalServerError:
		metrics.Collector.IncrementCounter(metrics.PusherStatus50x)
	}

	if response.StatusCode != http.StatusOK {
		return server.NewError(response.StatusCode, response.String())
	}

	return nil
}
```

### Metrics

Explanation

```
avg by(app, env, scope) (rate(pusher_success[$__rate_interval]))
avg by(app, env, scope) (rate(pusher_error[$__rate_interval]))
avg by(app, env, scope) (rate(pusher_http_20x[$__rate_interval]))
avg by(app, env, scope) (rate(pusher_http_40x[$__rate_interval]))
avg by(app, env, scope) (rate(pusher_http_50x[$__rate_interval]))
avg by(app, env, scope) (rate(pusher_http_timeoutx[$__rate_interval]))
```

#### Pusher dashboard

TODO Pusher Success, Pusher Errors, HTTP Time, 20x, 40x, 50x

<img width="1267" alt="image" src="https://user-images.githubusercontent.com/760657/221233377-c0dcc50b-8fa8-4064-9883-b0f7d05bc3fe.png">

Example

### Contributors

Fork me.

### Support

ariel.pineiro@iskaypet.com
