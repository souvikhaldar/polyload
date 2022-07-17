package uploader

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
)

type LocalStorage struct{}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (l *LocalStorage) UploadFile(fileName string, data []byte) error {
	dst, err := os.Create(
		viper.GetString("upload_dir") + "/" + fileName,
	)
	if err != nil {
		return errors.New(
			"Could not create the upload destination file: " + err.Error(),
		)
	}
	defer dst.Close()

	_, err = dst.Write(data)
	if err != nil {
		return errors.New(
			"Unable to copy files: " + err.Error(),
		)
	}
	return nil
}

func (l *LocalStorage) DownloadFile(fileName string) ([]byte, error) {
	log.Println("Downloading file: ", viper.GetString("upload_dir")+"/"+fileName)
	f, err := os.Open(
		viper.GetString("upload_dir") + "/" + fileName,
	)
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(f)
	if len(fileBytes) == 0 {
		return nil, errors.New("Empty file read")
	}
	return fileBytes, err
}
