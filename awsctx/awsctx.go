package awsctx

import (
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type AWSContext struct {
	SNSSvc snsiface.SNSAPI
}
