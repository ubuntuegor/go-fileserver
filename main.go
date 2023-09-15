package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const dirPerm = 0755

// differs from http.ServeFile in that there are no special redirects
func serveFile(w http.ResponseWriter, r *http.Request, path string) {
	file, err := os.Open(path)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Could not open file", http.StatusInternalServerError)
		}
		return
	}

	defer file.Close()

	modTime := time.UnixMilli(0)
	stat, err := os.Stat(path)

	if err == nil {
		modTime = stat.ModTime()
	}

	fileName := filepath.Base(path)
	http.ServeContent(w, r, fileName, modTime, file)
}

func writeFile(w http.ResponseWriter, r *http.Request, path string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, dirPerm)

	if err != nil {
		http.Error(w, "Could not create requested directories", http.StatusInternalServerError)
		return
	}

	file, err := os.Create(path)

	if err != nil {
		http.Error(w, "Could not create or overwrite file", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(file, r.Body)

	if err != nil {
		http.Error(w, "Could not write to file", http.StatusInternalServerError)
	}

	fmt.Fprint(w, "ok")
}

func deleteFile(w http.ResponseWriter, path string) {
	err := os.Remove(path)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Could not delete file", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprint(w, "ok")
}

func getRequestHandler(serveFromFolder string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path

		if strings.Contains(reqPath, "..") || strings.HasSuffix(reqPath, "/") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(serveFromFolder, reqPath)

		fmt.Printf("Got %v request for: %v\n", r.Method, filePath)

		switch r.Method {
		case http.MethodGet:
			serveFile(w, r, filePath)
		case http.MethodPost:
			writeFile(w, r, filePath)
		case http.MethodDelete:
			deleteFile(w, filePath)
		default:
			http.Error(w, "Unknown method", http.StatusBadRequest)
		}
	}
}

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		fmt.Println("Port is not specified or specified incorrectly in PORT environment variable")
		return
	}

	serveFromFolder := os.Getenv("SERVE_FROM_FOLDER")

	http.HandleFunc("/", getRequestHandler(serveFromFolder))

	fmt.Printf("Listening on port %v, serving folder %v\n", port, serveFromFolder)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
