package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ThumbnailError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeResponseError(w http.ResponseWriter, errMsg string, httpCode int) {
	js, err := json.Marshal(ThumbnailError{Code: httpCode, Message: errMsg})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(js)
	if err != nil {
		log.Fatal("Cannot send response!")
	}
}

func downloadImage(url string) ([]byte, error) {
	response, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	if response.StatusCode != 200 {
		err := errors.New(
			fmt.Sprintf("bad response during fetching image, status code %v", response.StatusCode))
		return nil, err
	}

	defer response.Body.Close()
	rawData, e := ioutil.ReadAll(response.Body)

	if e != nil {
		return nil, e
	}

	return rawData, nil
}
