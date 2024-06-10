package awsSdk

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitializeClient() (client *s3.Client, err error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return
	}
	client = s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(GetBaseEndpoint())
	})
	return
}

func GetBaseEndpoint() string {
	slices := []string{"https://", os.Getenv("AWS_REGION"), ".", os.Getenv("AWS_ENDPOINT_URL")}
	return strings.Join(slices, "")
}
