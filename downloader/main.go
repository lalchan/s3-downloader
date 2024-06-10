package downloader

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/schollz/progressbar/v3"

	"io"
	"os"
	"path/filepath"
)

func SaveObjects(client *s3.Client, files []string) (err error) {
	bar := progressbar.Default(int64(len(files)))

	for _, file := range files {
		err = SaveObject(client, bar, file)
		if err != nil {
			return
		}
	}
	return
}

func SaveObject(client *s3.Client, bar *progressbar.ProgressBar, filePath string) (err error) {
	input := s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(filePath),
	}
	object, err := client.GetObject(context.TODO(), &input)
	if err != nil {
		return err
	}

	go AsyncSaveFile(object, bar, filePath)
	return nil
}

func AsyncSaveFile(object *s3.GetObjectOutput, bar *progressbar.ProgressBar, path string) {
	bytes, err := io.ReadAll(object.Body)
	if err != nil {
		return
	}
	directory := filepath.Join(os.Getenv("SAVING_DIRECTORY"), filepath.Dir(path))
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return
	}
	fullPath := filepath.Join(directory, filepath.Base(path))
	file, err := os.Create(fullPath)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		return
	}
	bar.Add64(1)
}
