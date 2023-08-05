package utils

import (
	"io"
	"os"

	"github.com/disintegration/imaging"
)

func ResizeImage(reader NamedFileReader, width int, height int) (*os.File, error) {
	fileExt := GetFileExtension(reader.Name())
	tmpFile, err := os.CreateTemp(os.TempDir(), "dp-*."+fileExt)
	if err != nil {
		return nil, err
	}

	reader.Seek(0, io.SeekStart)
	_, err = io.Copy(tmpFile, reader)
	if err != nil {
		return nil, err
	}

	image, err := imaging.Open(tmpFile.Name(), imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	dst := imaging.Fill(image, width, height, imaging.Center, imaging.Lanczos)
	return tmpFile, imaging.Save(dst, tmpFile.Name())
}
