package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/d-smith/push/awsctx"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/ses"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
)

type RegisterSender struct {
	SenderEmail string `json:"sender_email"`
}

func processRegisterSender(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("processRegisterSender")
	var sender RegisterSender
	err := json.Unmarshal([]byte(request.Body), &sender)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	fmt.Println("register sender", sender.SenderEmail)
	verifyEmailIdentityIn := ses.VerifyEmailIdentityInput {
		EmailAddress: aws.String(sender.SenderEmail),
	}

	_, err = awsContext.SESSvc.VerifyEmailIdentity(&verifyEmailIdentityIn)
	if err != nil {
		fmt.Println("Error registering sender", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
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
	svc := ses.New(sess)

	awsContext.SESSvc = svc
	handler := makeHandler(&awsContext)
	lambda.Start(handler)
}

