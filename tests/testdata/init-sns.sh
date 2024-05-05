#!/bin/sh

echo "Initializing SNS topics..."

awslocal sns create-topic \
    --name OrderProductionTopic

awslocal sns create-topic \
    --name UpdateOrderTopic