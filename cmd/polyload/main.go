package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/souvikhaldar/polyload/pkg/uploader"
	"github.com/spf13/viper"
)

type server struct {
	load   uploader.Loader
	router *mux.Router
}

func init() {
	// parse config file
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.polyload") // call multiple times to add many search paths

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()
	viper.WatchConfig()
	// create the dir to upload files

	if err := os.MkdirAll(viper.GetString("upload_dir"), os.ModePerm); err != nil {
		log.Println("Unable to create upload dir "+viper.GetString("upload_dir")+": ", err)
	}
	if err := os.MkdirAll(viper.GetString("download_dir"), os.ModePerm); err != nil {
		log.Println("Unable to create download dir "+viper.GetString("upload_dir")+": ", err)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) isRegisteredUser(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO: this logic for verifying registered users needs to be updated to standard practices
		token := r.URL.Query().Get("token")

		readFile, err := os.Open(viper.GetString("registered_users"))
		if err != nil {
			log.Println("Could not read register user file: ", err)
			http.NotFound(w, r)
			return
		}
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)
		var fileLines []string

		for fileScanner.Scan() {
			fileLines = append(fileLines, fileScanner.Text())
		}

		readFile.Close()

		verified := false
		for _, line := range fileLines {
			if token == line {
				verified = true
			}
		}
		if !verified {
			log.Println("Malicious request from: ", r.RemoteAddr)
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
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

		// read the file contents
		fileData, err := ioutil.ReadAll(file)
		if len(fileData) == 0 {
			http.Error(
				w,
				"File size receiver is 0",
				http.StatusInternalServerError,
			)
			return

		}
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// which cloud provider to upload to
		cloud := r.URL.Query().Get("cloud")
		log.Println("Upload to cloud: ", cloud)

		switch cloud {
		case "azure":
			s.load = uploader.NewAzure()
		case "local":
			s.load = uploader.NewLocalStorage()
		case "aws":
			s.load = uploader.NewAws()
		default:
			s.load = uploader.NewLocalStorage()
		}
		if err := s.load.UploadFile(fileHeader.Filename, fileData); err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		w.WriteHeader(http.StatusOK)

		log.Println("Upload success")
		return
	}

}

func (s *server) handleFileDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running in handleFileDownload")
		// which cloud provider to upload to
		cloud := r.URL.Query().Get("cloud")
		fileName := r.URL.Query().Get("file")
		log.Println("Upload to cloud: ", cloud)
		log.Println("file: ", fileName)

		switch cloud {
		case "azure":
			s.load = uploader.NewAzure()
		case "local":
			s.load = uploader.NewLocalStorage()
		}
		fileBytes, err := s.load.DownloadFile(fileName)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return

		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
		return

	}
}

func (s *server) setRoutes() {
	s.router.HandleFunc(
		"/upload",
		s.isRegisteredUser(s.handleFileUpload()),
	).Methods("POST")

	s.router.HandleFunc(
		"/download",
		s.isRegisteredUser(s.handleFileDownload()),
	).Methods("GET")
}
func main() {
	// create a new server
	s := &server{
		uploader.NewAws(),
		mux.NewRouter(),
	}
	s.setRoutes()
	srv := http.Server{
		Addr:        ":" + viper.GetString("port"),
		Handler:     s,
		ReadTimeout: 0,
	}
	log.Fatal(srv.ListenAndServe())

}
