package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Metadata struct {
	Url string `json:"url"`
}

func main() {
	apiKey := os.Getenv("CAT_API_KEY")
	if apiKey == "" {
		log.Fatal("CAT_API_KEY is not set")
	}

	// get cat image
	getBreeds, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search?breed_ids=mcoo", nil)
	if err != nil {
		log.Fatalf("unable to create cat image request: %v", err)
	}
	getBreeds.Header.Set("x-api-key", apiKey)

	// grab cat metadata
	c := new(http.Client)
	breedsResp, err := c.Do(getBreeds)
	if err != nil {
		log.Fatalf("unable to get cat metadata: %v", err)
	}
	defer breedsResp.Body.Close()
	if breedsResp.StatusCode != http.StatusOK {
		log.Fatalf("received non 200 status code: %d", breedsResp.StatusCode)
	}
	var breeds []Metadata
	if err := json.NewDecoder(breedsResp.Body).Decode(&breeds); err != nil {
		log.Fatalf("unable to unmarshal cat metadata: %v", err)
	}

	if len(breeds) == 0 {
		log.Fatal("no cat metadata found")
	}
	fmt.Println("found", len(breeds)+1, "breeds")

	// download image
	getImage, err := http.NewRequest("GET", breeds[0].Url, nil)
	if err != nil {
		log.Fatalf("unable to create cat image request: %v", err)
	}
	getImage.Header.Set("x-api-key", apiKey)

	imageResp, err := c.Do(getImage)
	if err != nil {
		log.Fatalf("unable to get cat image: %v", err)
	}
	defer imageResp.Body.Close()
	if imageResp.StatusCode != http.StatusOK {
		log.Fatalf("received non 200 status code: %d", imageResp.StatusCode)
	}
	fmt.Println("retrieving image")

	// upload image to s3
	sess, err := session.NewSession(aws.NewConfig().WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to create aws session: %v", err)
	}

	key := "cat-" + strconv.Itoa(time.Now().Minute()) + ".jpg"
	uploader := s3manager.NewUploader(sess)
	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("go-kt"),
		Key:    &key,
		Body:   imageResp.Body,
	}); err != nil {
		log.Fatalf("unable to upload cat image: %v", err)
	}
	fmt.Println("uploaded image to s3")

	f, err := os.Create(key)
	if err != nil {
		log.Fatalf("unable to create file: %v", err)
	}
	downloader := s3manager.NewDownloader(sess)
	if _, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String("go-kt"),
		Key:    &key,
	}); err != nil {
		log.Fatalf("unable to download cat image: %v", err)
	}

	fmt.Println("downloaded from S3")
}
