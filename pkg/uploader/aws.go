package uploader

type Aws struct {
}

func NewAws() *Aws {
	return &Aws{}
}
func (a *Aws) UploadFile(fileName string, data []byte) error {
	return nil
}
func (a *Aws) DownloadFile(fileName string) ([]byte, error) {
	return nil, nil
}
