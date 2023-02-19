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

TODO

### Consumer

TODO

### Pusher

TODO

### RestClient

TODO

#### RestClient configuration

TODO

#### RestClient usage

TODO

### Metrics

Explanation

#### Pusher dashboard

TODO Pusher Success, Pusher Errors, HTTP Time, 20x, 40x, 50x

### Contributors

TODO

### Support

TODO


