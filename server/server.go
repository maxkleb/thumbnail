package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/maxkleb/thumbnail/transformation"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ThumbnailError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.Println("Starting thumbnail server with port "+ port +"...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/thumbnail", thumbnailHandler).Methods("GET")
	address := ":" + port
	log.Fatal(http.ListenAndServe(address, router))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("Health check has called")
	w.WriteHeader(http.StatusOK)
}

func writeResponseError(w http.ResponseWriter, errMsg string, httpCode int) {
	js, err := json.Marshal(ThumbnailError{Code: httpCode, Message: errMsg})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}


func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil {
		writeResponseError(w, "invalid input for width", http.StatusBadRequest)
		return
	}

	y, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil {
		writeResponseError(w, "invalid input for height", http.StatusBadRequest)
		return
	}

	url := r.URL.Query().Get("url")

	imgBytes, err := downloadImage(url)
	if err != nil {
		writeResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if imgBytes == nil {
		writeResponseError(w, "cannot open provided image", http.StatusBadRequest)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		writeResponseError(w, "error during image decoding ", http.StatusBadRequest)
		return
	}

	finalImage := transformation.ProcessImg(x, y, img)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, finalImage, nil)
	resBytes := buf.Bytes()

	// Encode the image data and write as an server response
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(resBytes)))
	if _, err = w.Write(resBytes); err != nil {
		writeResponseError(w,"unable to reconstruct image", http.StatusInternalServerError)
		return
	}

}