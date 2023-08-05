package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Image struct {
	Width    uint   `json:"width"`
	Height   uint   `json:"height"`
	Filename string `json:"image,omitempty"`
	URL      string `json:"url,omitempty"`
}

type Images []*Image

func (i *Images) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if bytes == nil {
		return nil
	}

	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &i)
}

func (i Images) Value() (driver.Value, error) {
	if len(i) == 0 {
		return nil, nil
	}

	return json.Marshal(i)
}

func (i *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if bytes == nil {
		return nil
	}

	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &i)
}

func (image *Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	newStruct := &struct {
		*Alias
	}{
		Alias: (*Alias)(image),
	}

	if len(newStruct.URL) > 0 {
		newStruct.Filename = ""
	}

	return json.Marshal(newStruct)
}
