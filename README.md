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

## Table of contents

* [Project setup](#project-setup)
* [SDK](#sdk)
  * [Configuration](#configuration)
  * [Consumer](#consumer)
  * [Pusher](#Pusher)
  * [RestClient](#restclient)
    * [RestClient configuration](#restclient-configuration)
    * [RestClient usage](#restclient-usage)
  * [Metrics](#metrics)
    * [Pusher dashboard](#tasks-dashboard)
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

```yaml
# aws
aws:
  url: http://localhost:4566
  id: development
  secret: development
  profile: default
  region: us-east-1
```

### Consumer

Queue to consume messages.

```yaml
# consumers
consumers:
  users:
    queue-url: http://localhost:4566/000000000000/my-queue #your-amazon-sqs-queue
```

### Pusher

Your app to receive messages. Example: my.app/news. Must allow POST Http Request in
Application/json.

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
  endpoint: http://localhost:8000/my-post-endpoint #
```

### RestClient

Pusher app need a rest client to send messages to target.

#### RestClient configuration

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

#### RestClient usage

```go
func (c HttpAppClient) PostMessage(requestBody *RequestBody) error {
	startTime := time.Now()
	response := c.rb.Post(c.baseURL, requestBody)
	elapsedTime := time.Since(startTime)

	metrics.Collector.RecordExecutionTime("consumers.pusher.http.time", elapsedTime.Milliseconds())

	if response.Err != nil {
		return response.Err
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		metrics.Collector.IncrementCounter("consumers.pusher.http.20x")
	} else {
		if response.StatusCode >= 400 && response.StatusCode < 500 {
			metrics.Collector.IncrementCounter("consumers.pusher.http.40x")
		} else {
			if response.StatusCode >= 500 {
				metrics.Collector.IncrementCounter("consumers.pusher.http.50x")
			}
		}
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
consumers.pusher.http.20x: messages that were sent and confirmed successfully
consumers.pusher.http.40x: messages that weren't sent or confirmed successfully
consumers.pusher.http.50x: messages that weren't sent or confirmed successfully
consumers_pusher_http_time: delivery time, remember configure your rest client correctly
```

#### Pusher dashboard

TODO Pusher Success, Pusher Errors, HTTP Time, 20x, 40x, 50x

### Contributors

Fork me

### Support

arielsrv@gmail.com


