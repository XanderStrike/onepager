package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/peterbourgon/diskv"
)

type HomePage struct {
	Files []os.FileInfo
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./pages/")
	if err != nil {
		log.Fatal(err)
	}
	render(w, "index.html", HomePage{
		Files: files,
	})
}

func flatTransform(s string) []string { return []string{} }
func NewHandler(w http.ResponseWriter, r *http.Request) {
	filename := fmt.Sprintf("%s.html", r.FormValue("filename"))
	content := r.FormValue("content")
	d := diskv.New(diskv.Options{
		BasePath:     "pages",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	d.Write(filename, []byte(content))
	io.WriteString(w, fmt.Sprintf("filename: %s\ncontents: %s", filename, content))
	log.Println("Saved", filename)
}

func main() {
	log.Println("Starting!")
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/new", NewHandler)
	r.PathPrefix("/pages").Handler(http.FileServer(http.Dir("./")))
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
