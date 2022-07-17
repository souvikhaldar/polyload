package uploader

import (
	"bytes"
	"context"
	"errors"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/viper"
)

type Azure struct {
	ctx             context.Context
	url             string
	serviceClient   *azblob.ServiceClient
	containerClient *azblob.ContainerClient
	credential      *azidentity.DefaultAzureCredential
}

func NewAzure() *Azure {
	ctx := context.Background()
	az := new(Azure)
	az.ctx = ctx
	url := "https://" + viper.GetString("azure.blob_storage_account_name") + ".blob.core.windows.net/"
	az.url = url
	// Create a default Azure credential
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	az.credential = credential
	serviceClient, err := azblob.NewServiceClient(url, credential, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	az.serviceClient = serviceClient
	containerClient, err := serviceClient.NewContainerClient(viper.GetString("azure.blob_container_name"))
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	az.containerClient = containerClient

	return az
}

func (a *Azure) UploadFile(fileName string, data []byte) error {
	blobClient, err := azblob.NewBlockBlobClient(
		a.url+viper.GetString("azure.blob_container_name")+"/"+fileName,
		a.credential,
		nil,
	)
	if err != nil {
		return errors.New("Couldn't create a blog client: " + err.Error())
	}
	// Upload to data to blob storage
	uploadOption := azblob.UploadOption{}
	_, err = blobClient.UploadBuffer(
		a.ctx,
		data,
		uploadOption,
	)

	if err != nil {
		return errors.New("Failure to upload to blob: " + err.Error())
	}
	return nil
}

func (a *Azure) DownloadFile(fileName string) ([]byte, error) {
	blobClient, err := azblob.NewBlockBlobClient(
		a.url+viper.GetString("azure.blob_container_name")+"/"+fileName,
		a.credential,
		nil,
	)
	if err != nil {
		return nil, errors.New("Couldn't create a blog client: " + err.Error())
	}
	// Download the blob
	get, err := blobClient.Download(a.ctx, nil)
	if err != nil {
		return nil, errors.New("Could not download " + err.Error())
	}
	downloadedData := &bytes.Buffer{}
	reader := get.Body(&azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return nil, errors.New("Could not read " + err.Error())
	}
	err = reader.Close()
	if err != nil {
		return nil, errors.New("Could not close " + err.Error())
	}

	return downloadedData.Bytes(), nil
}
