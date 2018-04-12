package awsctx

import (
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

type AWSContext struct {
	SNSSvc snsiface.SNSAPI
	SESSvc sesiface.SESAPI
}
