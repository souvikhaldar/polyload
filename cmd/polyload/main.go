package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/souvikhaldar/polyload/pkg/uploader"
	"github.com/spf13/viper"
)

type server struct {
	upload uploader.Uploader
	router *mux.Router
}

func init() {
	// parse config file
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()
	viper.WatchConfig()
	// create the dir to upload files

	if err := os.MkdirAll(viper.GetString("upload_dir"), os.ModePerm); err != nil {
		log.Println("Unable to create upload dir "+viper.GetString("upload_dir")+": ", err)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleFileUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running in file upload")

		// process uploaded files of max 50 mb
		r.Body = http.MaxBytesReader(
			w,
			r.Body,
			viper.GetInt64("max_upload_size"),
		)
		if err := r.ParseMultipartForm(
			viper.GetInt64("max_upload_size"),
		); err != nil {
			http.Error(w, "Uploaded file too big", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		log.Println("File name: ", fileHeader.Filename)
		dst, err := os.Create(
			viper.GetString("upload_dir") + "/" + fileHeader.Filename,
		)
		if err != nil {
			http.Error(
				w,
				"Could not create the upload destination file: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(
				w,
				"Unable to copy files: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		log.Println("Upload success")
		return
	}

}

func (s *server) setRoutes() {
	s.router.HandleFunc("/upload", s.handleFileUpload()).Methods("POST")
}
func main() {
	// create a new server
	s := &server{
		uploader.NewAws(),
		mux.NewRouter(),
	}
	s.setRoutes()
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), s))

}
