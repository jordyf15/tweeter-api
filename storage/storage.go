package storage

import (
	"io"
	"sync"

	"github.com/jordyf15/tweeter-api/models"
)

type Storage interface {
	UploadFile(respond chan<- error, wg *sync.WaitGroup, file io.ReadSeeker, key string, metadata map[string]string)
	RemoveFile(respond chan<- error, wg *sync.WaitGroup, key string)
	GetFileLink(key string) (string, error)
	AssignImageURLToUser(model *models.User)
	AssignImageURLToGroup(model *models.Group)
}
