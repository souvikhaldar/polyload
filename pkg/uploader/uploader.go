package uploader

type Uploader interface {
	UploadFile(fileName string) (int64, error)
	DownloadFile(fileNumber int64) ([]byte, error)
}
