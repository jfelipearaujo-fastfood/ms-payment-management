#!/bin/sh

echo "Initializing SQS queues..."

awslocal sqs create-queue \
    --queue-name OrderPaymentQueue