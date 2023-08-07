package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"github.com/jordyf15/tweeter-api/models"
	"google.golang.org/api/option"
)

type cloudStorage struct {
	client *storage.Client
	ctx    context.Context
}

func NewCloudStorage() Storage {
	config := &firebase.Config{
		StorageBucket: os.Getenv("GCP_BUCKET_NAME"),
	}
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	fbClient, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return &cloudStorage{client: fbClient, ctx: ctx}
}

func (api *cloudStorage) UploadFile(respond chan<- error, wg *sync.WaitGroup, file io.ReadSeeker, key string, metadata map[string]string) {
	if wg != nil {
		defer wg.Done()
	}

	bucket, err := api.client.Bucket(os.Getenv("GCP_BUCKET_NAME"))
	if err != nil {
		respond <- err
		return
	}

	obj := bucket.Object(key)
	ctx, cancel := context.WithTimeout(api.ctx, time.Second*30)
	defer cancel()

	wc := obj.NewWriter(ctx)
	wc.Metadata = metadata

	file.Seek(0, 0)
	if _, err := io.Copy(wc, file); err != nil {
		respond <- err
		return
	}
	if err := wc.Close(); err != nil {
		respond <- err
		return
	}

	respond <- nil
}

func (api *cloudStorage) GetFileLink(key string) (string, error) {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", os.Getenv("GCP_BUCKET_NAME"), key), nil
}

func (storage *cloudStorage) AssignImageURLToUser(user *models.User) {
	for _, img := range user.ProfileImages {
		img.URL, _ = storage.GetFileLink(user.ImagePath(img))
	}

	user.BackgroundImage.URL, _ = storage.GetFileLink(user.ImagePath(&user.BackgroundImage))
}

func (api *cloudStorage) RemoveFile(respond chan<- error, wg *sync.WaitGroup, key string) {
	if wg != nil {
		defer wg.Done()
	}

	bucket, err := api.client.Bucket(os.Getenv("GCP_BUCKET_NAME"))
	if err != nil {
		respond <- err
		return
	}

	ctx, cancel := context.WithTimeout(api.ctx, time.Second*30)
	defer cancel()

	err = bucket.Object(key).Delete(ctx)
	if err != nil {
		fmt.Printf("gcp error: %v for key %s\n", err, key)
	}

	respond <- err
}
