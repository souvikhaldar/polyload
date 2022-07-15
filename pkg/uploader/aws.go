package uploader

type Aws struct {
}

func NewAws() *Aws {
	return &Aws{}
}
func (a *Aws) UploadFile(fileName string) (int64, error) {
	return 0, nil
}
func (a *Aws) DownloadFile(fileNumber int64) ([]byte, error) {
	return nil, nil
}
