#!/bin/bash
url=http://localhost:4566
region=us-east-1
topic=orders-topic
queue=orders-consumer
account=000000000000

awslocal sqs create-queue --queue-name $queue
awslocal sqs purge-queue --queue-url $url/$account/$queue
awslocal sqs list-queues
awslocal sqs receive-message --queue-url $url/$account/$queue
awslocal sns create-topic --name $topic
awslocal sns list-subscriptions
awslocal sns subscribe --topic-arn arn:aws:sns:$region:$account:$topic --protocol sqs --notification-endpoint arn:aws:sns:$region:$account:$queue
awslocal sqs get-queue-attributes --queue-url http://localhost:4566/$account/$queue --attribute-names All
awslocal sns publish --topic-arn arn:aws:sns:$region:$account:$topic --message '{"order_id": 1}'
awslocal secretsmanager create-secret --name "cache.password" --secret-string 'a300p011'
