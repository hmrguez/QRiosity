#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e


# Run the build command
cd src/frontend
bun run build

# Define the S3 bucket name
S3_BUCKET="qriosity-site"

# Upload the build output to the S3 bucket
aws s3 sync ./dist s3://$S3_BUCKET --delete

echo "Deployment to S3 bucket $S3_BUCKET completed successfully."