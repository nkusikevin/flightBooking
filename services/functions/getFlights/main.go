package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB
var sess *session.Session

func init() {
	// Create a new session, allowing SDK to use the default credential chain
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new DynamoDB client
	svc = dynamodb.New(sess)
}

type Flight struct {
	PK            string    `json:"PK"`
	SK            string    `json:"SK"`
	FlightNumber  string    `json:"flightNumber"`
	Airline       string    `json:"airline"`
	Price         float64   `json:"price"`
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DepartureDate string    `json:"departureDate"`
	DepartureTime string    `json:"departureTime"`
	ArrivalDate   string    `json:"arrivalDate"`
	ArrivalTime   string    `json:"arrivalTime"`
	ClientNames   []string  `json:"clientNames"`
	CreatedAt     time.Time `dynamodbav:"CreatedAt"`
	UpdatedAt     time.Time `dynamodbav:"UpdatedAt"`
}

func getFlights(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create the input for the query
	input := &dynamodb.ScanInput{
		TableName: aws.String("flights"),
	}

	// Retrieve the items from the DynamoDB table
	result, err := svc.Scan(input)

	if err != nil {
		log.Fatalf("Got error calling Scan: %s", err)
		return events.APIGatewayProxyResponse{Body: "Error getting flights", StatusCode: 500}, nil
	}

	// Initialize the Flights array
	flights := []Flight{}

	// Unmarshal the Items field in the result value to the Flights array
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &flights)
	if err != nil {
		log.Fatalf("Failed to unmarshal Scan Items, %v", err)
		return events.APIGatewayProxyResponse{Body: "Error getting flights", StatusCode: 500}, nil
	}

	// Marshal the Flights array into JSON
	body, err := json.Marshal(flights)
	if err != nil {
		log.Fatalf("Failed to marshal flights, %v", err)
		return events.APIGatewayProxyResponse{Body: "Error getting flights", StatusCode: 500}, nil
	}

	// Return a response with a 200 OK status and the JSON flights
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	lambda.Start(getFlights)
}
