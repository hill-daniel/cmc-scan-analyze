package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

func main() {
	lambda.Start(handler)
}
