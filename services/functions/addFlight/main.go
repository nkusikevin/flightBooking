package main

import (
	"encoding/json"
	"log"

	// "net/http"

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
	Origin        string    `json:"origin"`
	Price         float64   `json:"price"`
	Destination   string    `json:"destination"`
	DepartureDate string    `json:"departureDate"`
	DepartureTime string    `json:"departureTime"`
	ArrivalDate   string    `json:"arrivalDate"`
	ArrivalTime   string    `json:"arrivalTime"`
	ClientNames   string    `json:"clientNames"`
	CreatedAt     time.Time `dynamodbav:"CreatedAt"`
	UpdatedAt     time.Time `dynamodbav:"UpdatedAt"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//unmarshal the request body
	var flight Flight
	err := json.Unmarshal([]byte(request.Body), &flight)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error unmarshalling request", StatusCode: 500}, nil
	}

	// Create a new Flight item
	flightItem := Flight{
		PK:            "FLIGHT#" + flight.FlightNumber,
		SK:            "METADATA#" + flight.FlightNumber,
		FlightNumber:  flight.FlightNumber,
		Price:         flight.Price,
		Origin:        flight.Origin,
		Destination:   flight.Destination,
		DepartureDate: flight.DepartureDate,
		DepartureTime: flight.DepartureTime,
		ArrivalDate:   flight.ArrivalDate,
		ArrivalTime:   flight.ArrivalTime,
		ClientNames:   flight.ClientNames,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Marshal the Flight struct into an AWS SDK-compatible type
	av, err := dynamodbattribute.MarshalMap(flightItem)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
		return events.APIGatewayProxyResponse{Body: "Error marshalling map", StatusCode: 500}, nil
	}

	// Create a new item in the table
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("flights"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return events.APIGatewayProxyResponse{Body: "Error calling PutItem", StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: "Flight created", StatusCode: 201}, nil

}

func main() {
	lambda.Start(handler)
}
