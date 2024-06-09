package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/briandowns/spinner"
)

type GetObjectsReturnType struct {
	Files             []string
	ContinuationToken *string
}

func main() {
	client, err := initializeClient()
	if err != nil {
		log.Fatal(err)
	}
	fetchCompleted := false
	var continuationToken *string
	var files []string
	for !fetchCompleted {
		objects, err := GetObjectList(client, continuationToken)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, objects.Files...)
		continuationToken = objects.ContinuationToken
		fetchCompleted = (objects.ContinuationToken == nil)
	}
	log.Printf("Total count: %d", len(files))

	spinner := spinner.New(spinner.CharSets[41], 100*time.Millisecond)
	spinner.Suffix = "Saving files..."
	spinner.Start()
	defer spinner.Stop()
	time.Sleep(2 * time.Second)

	for _, file := range files {
		err = SaveObject(client, file)
		if err != nil {
			log.Fatal(err)
		}
	}
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

func SaveObject(client *s3.Client, filePath string) (err error) {
	input := s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(filePath),
	}
	object, err := client.GetObject(context.TODO(), &input)
	if err != nil {
		return err
	}

	go AsyncSaveFile(object, filePath)
	return nil
}

func AsyncSaveFile(object *s3.GetObjectOutput, path string) {
	bytes, err := io.ReadAll(object.Body)
	if err != nil {
		return
	}
	fullPath := filepath.Join(os.Getenv("SAVING_DIRECTORY"), path)
	file, err := os.Create(fullPath)

	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(bytes)
	if err != nil {
		return
	}
}

func initializeClient() (client *s3.Client, err error) {
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
