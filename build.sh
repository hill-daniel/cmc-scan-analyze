#bin/bash
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap && zip function.zip bootstrap
rm bootstrap
