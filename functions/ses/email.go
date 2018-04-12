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
	EmailAddress string `json:"email_address"`
}

type SendEmailSpec struct {
	Sender string `json:"sender_email"`
	To []string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type SESSendResponse struct {
	MessageID string `json:"messageId"`
}

func processRegisterSender(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("processRegisterSender")
	var sender RegisterSender
	err := json.Unmarshal([]byte(request.Body), &sender)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	fmt.Println("register sender", sender.EmailAddress)
	verifyEmailIdentityIn := ses.VerifyEmailIdentityInput {
		EmailAddress: aws.String(sender.EmailAddress),
	}

	_, err = awsContext.SESSvc.VerifyEmailIdentity(&verifyEmailIdentityIn)
	if err != nil {
		fmt.Println("Error registering sender", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func ssToDestination(ss []string) *ses.Destination {
	var dest ses.Destination
	for _,s := range ss {
		dest.ToAddresses = append(dest.ToAddresses, aws.String(s))
	}

	return &dest
}

func processSendEmail(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(processSendEmail)

	var mailMessage SendEmailSpec
	err := json.Unmarshal([]byte(request.Body), &mailMessage)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}


	fmt.Printf("send message %+v\n", mailMessage)
	sendEmailInput := ses.SendEmailInput{
		Destination: ssToDestination(mailMessage.To),
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(mailMessage.Subject),
			},
			Body:&ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(mailMessage.Message),
				},
			},
		},
		Source: aws.String(mailMessage.Sender),
	}

	out, err := awsContext.SESSvc.SendEmail(&sendEmailInput)
	if err != nil {
		fmt.Println("Error sending email", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	responseJson := SESSendResponse{
		MessageID: *out.MessageId,
	}

	responseOut, _ := json.Marshal(&responseJson)

	return events.APIGatewayProxyResponse{Body: string(responseOut),StatusCode: 200}, nil
}

func processRequest(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("processing email request", request.Path)

	if request.Path == "/push/api/v1/regaddress" {
		return processRegisterSender(awsContext, request)
	}

	return processSendEmail(awsContext, request)
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

