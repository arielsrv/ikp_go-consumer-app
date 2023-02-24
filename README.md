> This SDK provides a framework to build middle-end for consumers and APIs

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

<img width="1211" alt="image" src="https://user-images.githubusercontent.com/760657/220067058-3f86739d-eaa5-49fd-9000-663d99f81642.png">

## Table of contents

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

## Project setup

TODO specify, topic name (local), consumer name, target pusher,

```shell
awslocal sqs create-queue --queue-name orders-consumer
awslocal sqs purge-queue --queue-url http://localhost:4566/000000000000/orders-consumer
awslocal sqs list-queues
awslocal sqs receive-message --queue-url http://localhost:4566/000000000000/orders-consumer
awslocal sns create-topic --name orders-topic
awslocal sns list-subscriptions
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:orders-topic --protocol sqs --notification-endpoint arn:aws:sns:us-east-1:000000000000:orders-consumer
awslocal sns publish --topic-arn arn:aws:sns:us-east-1:000000000000:orders-topic --message '{"order_id": 1}'
```

## SDK

### Configuration

Environment configuration is based on **Archaius Config**, you should use a similar folder
structure.
*SCOPE* env variable in remote environment is required

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

* resources/config/config.properties (shared config)
* resources/config/{environment}/config.properties (override shared config by environment)
* resources/config/{environment}/{scope}.config.properties (override env and shared config by scope)

example *test.pets-api.internal.com*

```
└── config
    ├── config.yml                              3th (third)
    └── local
        └── config.yml                          <ignored>
    └── prod
        └── config.yml (base config)            2nd (second)
        └── test.config.yml (base config)       1st (first)
```

* 1st (first)   prod/test.config.yml
* 2nd (second)  prod/config.yml
* 3th (third)   config.yml

```
2022/11/20 13:24:26 INFO: Two files have same priority. keeping
    /resources/config/prod/test.config.yml value
2022/11/20 13:24:26 INFO: Configuration files:
    /resources/config/prod/test.config.yml,
    /resources/config/prod/config.yml,
    /resources/config/config.yml
2022/11/20 13:24:26 INFO: invoke dynamic handler:FileSource
2022/11/20 13:24:26 INFO: enable env source
2022/11/20 13:24:26 INFO: invoke dynamic handler:EnvironmentSource
2022/11/20 13:24:26 INFO: archaius init success
2022/11/20 13:24:26 INFO: ENV: prod, SCOPE: test
2022/11/20 13:24:26 INFO: create new watcher
2022/11/20 13:24:26 Listening on port 8080
2022/11/20 13:24:26 Open http://127.0.0.1:8080/ping in the browser
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

```go
func (c HttpAppClient) PostMessage(requestBody *RequestBody) error {
	startTime := time.Now()
	response := c.rb.Post(c.baseURL, requestBody)
	elapsedTime := time.Since(startTime)

	metrics.Collector.
		RecordExecutionTime("consumers.pusher.http.time", elapsedTime.Milliseconds())

	if response.Err != nil {
		return response.Err
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		metrics.Collector.IncrementCounter("consumers.pusher.http.20x")
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		metrics.Collector.IncrementCounter("consumers.pusher.http.40x")
	} else if response.StatusCode >= 500 {
		metrics.Collector.IncrementCounter("consumers.pusher.http.50x")
	}

	if response.StatusCode != http.StatusOK {
		return fiber.NewError(response.StatusCode, response.String())
	}

	return nil
}
```

### Metrics

Explanation

```
consumers_pusher_success: messages that were sent and confirmed successfully
consumers_pusher_errors: messages that weren't sent or confirmed successfully
consumers_pusher_http_20x: messages that were sent and confirmed successfully
consumers_pusher_http_40x: messages that weren't sent or confirmed successfully
consumers_pusher_http_50x: messages that weren't sent or confirmed successfully
consumers_pusher_http_time: delivery time, remember configure your rest client correctly
```

#### Pusher dashboard

TODO Pusher Success, Pusher Errors, HTTP Time, 20x, 40x, 50x

<img width="1267" alt="image" src="https://user-images.githubusercontent.com/760657/221233377-c0dcc50b-8fa8-4064-9883-b0f7d05bc3fe.png">

### Contributors

Fork me

### Support

arielsrv@gmail.com


