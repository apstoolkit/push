package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/d-smith/push/awsctx"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dynamoDBMockery struct {
	snsiface.SNSAPI
}

func (m *dynamoDBMockery) Publish(request *sns.PublishInput) (*sns.PublishOutput, error) {
	var out sns.PublishOutput

	out.MessageId = aws.String("m1")
	return &out, nil
}

func TestSMSPush(t *testing.T) {

	tests := []struct {
		name    string
		request events.APIGatewayProxyRequest
		expect  int
		err     error
	}{
		{
			"Handle parse error",
			events.APIGatewayProxyRequest{Body: `{"phoneNo:"`},
			400,
			nil,
		},
		{
			"Good request",
			events.APIGatewayProxyRequest{Body: `{"phoneNo":"1112223333", "message":"yeah"}`},
			200,
			nil,
		},
	}

	var awsContext awsctx.AWSContext
	var myMock dynamoDBMockery
	awsContext.SNSSvc = &myMock

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := processRequest(&awsContext, test.request)
			assert.IsType(t, test.err, err)
			assert.Equal(t, test.expect, response.StatusCode)
			if assert.NotNil(t, response) {
				fmt.Println("Body", response.Body)
			}
		})

	}
}
