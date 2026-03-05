#!/bin/bash

# Create bucket
awslocal s3 mb s3://ecommerce-uploads

# create queue
awslocal sqs create-queue --queue-name ecommerce-events

echo "LocalStack initialization complete"