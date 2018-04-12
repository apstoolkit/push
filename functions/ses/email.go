package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/d-smith/push/awsctx"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RegisterSender struct {
	SenderEmail string `json:"sender_email`
}

func processRegisterSender(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: string("register sender\n"), StatusCode: 200}, nil
}

func processRequest(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("processing email request", request.Path)

	if request.Path == "/push/api/v1/regsender" {
		return processRegisterSender(awsContext, request)
	}

	return events.APIGatewayProxyResponse{Body: string("got it"), StatusCode: 200}, nil
}

func makeHandler(awsContext *awsctx.AWSContext) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return processRequest(awsContext, request)
	}
}

func main() {

	var awsContext awsctx.AWSContext

	sess := session.New()
	svc := sns.New(sess)

	awsContext.SNSSvc = svc
	handler := makeHandler(&awsContext)
	lambda.Start(handler)
}

