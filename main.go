package main

import (
	"fmt"
	"log"

	"lalchan/download-s3-objects/awsSdk"
	"lalchan/download-s3-objects/downloader"
	"lalchan/download-s3-objects/fetcher"
)

var failedReading []string
var failedDiretory map[string]bool
var failedCreating []string
var failedWriting []string

type GetObjectsReturnType struct {
	Files             []string
	ContinuationToken *string
}

func main() {
	client, err := awsSdk.InitializeClient()
	if err != nil {
		log.Fatal(err)
	}

	files, err := fetcher.FetchAllFiles(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total count: %d", len(files))

	err = downloader.SaveObjects(client, files)
	if err != nil {
		log.Fatal(err)
	}
}

func printStringList(list []string) {
	for _, file := range list {
		fmt.Println(file)
	}
}

func printMap(m map[string]bool) {
	for k := range m {
		fmt.Println(k)
	}
}
