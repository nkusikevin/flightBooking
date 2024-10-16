package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ctx := context.Background()
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}

func handler(ctx context.Context, arguments *Input) (string, error) {
	fmt.Print("hello world")

	return arguments.Name, nil
}
