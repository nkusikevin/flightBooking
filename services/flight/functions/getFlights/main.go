package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ctx := context.Background()
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}

func handler(ctx context.Context, arguments *Input) (string, error) {
	fmt.Print("hello world")
	resp, err := json.MarshalIndent(arguments, "", "  ")
	if err != nil {
		return "", err
	}
	fmt.Print(string(resp))
	return arguments.Flight.Name, nil
}
