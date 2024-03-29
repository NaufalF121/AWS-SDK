package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	keyID := os.Getenv("AWS_ID")
	applicationKey := os.Getenv("AWS_KEY")
	bucketName := os.Getenv("BUCKET_NAME")
	endpoint := os.Getenv("ENDPOINT")
	region := os.Getenv("REGION")

	b2Client, err := NewClient(endpoint, keyID, applicationKey, bucketName, region, "")
	if err != nil {
		panic("Error in connecting to S3")
	}
	fmt.Println("Connected to S3")
	listResult, err := b2Client.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range listResult {
		log.Println(result)
	}

}
