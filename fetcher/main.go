package fetcher

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type GetObjectsReturnType struct {
	Files             []string
	ContinuationToken *string
}

func FetchAllFiles(client *s3.Client) (files []string, err error) {

	var continuationToken *string
	files = []string{}
	for true {
		object, loopErr := GetObjectList(client, continuationToken)
		if loopErr != nil {
			return files, loopErr
		}
		files = append(files, object.Files...)
		continuationToken = object.ContinuationToken
		if continuationToken == nil {
			break
		}
	}
	return
}

func GetObjectList(client *s3.Client, continuationToken *string) (output GetObjectsReturnType, err error) {
	output = GetObjectsReturnType{
		Files: []string{},
	}
	input := s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
	}
	if continuationToken != nil {
		input.ContinuationToken = continuationToken
	}
	objects, err := client.ListObjectsV2(context.TODO(), &input)
	if err != nil {
		return
	}
	output.ContinuationToken = objects.NextContinuationToken
	for _, object := range objects.Contents {
		output.Files = append(output.Files, *object.Key)
	}
	return
}
