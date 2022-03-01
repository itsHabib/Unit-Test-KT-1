package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go-kt-1/cats"
	"go-kt-1/storage"
)

func main() {
	apiKey := os.Getenv("CAT_API_KEY")
	if apiKey == "" {
		log.Fatal("CAT_API_KEY is not set")
	}

	catClient, err := cats.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewService("go-kt")
	if err != nil {
		log.Fatal(err)
	}

	breeds, err := catClient.GetBreeds("mcoo")
	if err != nil {
		log.Fatal(err)
	}

	if len(breeds) == 0 {
		log.Fatal("no cat metadata found")
	}
	fmt.Println("found", len(breeds)+1, "breeds")

	// download image from API
	image, err := catClient.GetCatImage(breeds[0].Url)
	if err != nil {
		log.Fatal(err)
	}

	key := "cat-" + strconv.Itoa(time.Now().Minute()) + ".jpg"
	if err := storage.Upload(image, key); err != nil {
		log.Fatal(err)
	}
	fmt.Println("uploaded image to s3")

	// download image to file
	if err := storage.Download(key, key); err != nil {
		log.Fatal(err)
	}

	fmt.Println("downloaded from S3")
}
