package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func processRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)

	return events.APIGatewayProxyResponse{Body: "hello\n", StatusCode: 200}, nil
}

func makeHandler() func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return processRequest(request)
	}
}

func main() {

	handler := makeHandler()
	lambda.Start(handler)
}
