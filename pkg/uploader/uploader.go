package uploader

type Loader interface {
	UploadFile(fileName string, data []byte) error
	DownloadFile(fileName string) ([]byte, error)
}
