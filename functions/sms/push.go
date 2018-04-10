package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/d-smith/push/awsctx"
)

type SMSMessage struct {
	PhoneNo string `json:"phoneNo"`
	Message string `json:"message"`
}

type SMSSendResponse struct {
	MessageID string `json:"messageId"`
}

var sender = &sns.MessageAttributeValue{
	DataType:    aws.String("String"),
	StringValue: aws.String("APSStatus"),
}

var messageAttrs = map[string]*sns.MessageAttributeValue{
	"AWS.SNS.SMS.SenderID": sender,
}

func processRequest(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)

	//Parse the body
	var msg SMSMessage

	err := json.Unmarshal([]byte(request.Body), &msg)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	publishInput := sns.PublishInput{
		PhoneNumber:       aws.String(msg.PhoneNo),
		Message:           aws.String(msg.Message),
		MessageAttributes: messageAttrs,
	}

	fmt.Println("Publishing message", *publishInput.Message, "to destination", *publishInput.PhoneNumber)

	out, err := awsContext.SNSSvc.Publish(&publishInput)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	responseBody := SMSSendResponse{
		MessageID: *out.MessageId,
	}

	responseOut, _ := json.Marshal(&responseBody)

	return events.APIGatewayProxyResponse{Body: string(responseOut), StatusCode: 200}, nil
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
