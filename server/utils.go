package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func downloadImage(url string) ([]byte, error) {
	response, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	if response.StatusCode != 200 {
		return nil, errors.New(
			fmt.Sprintf("bad response during fetching image, status code %v", response.StatusCode))
	}

	defer response.Body.Close()
	rawData, e := ioutil.ReadAll(response.Body)

	if e != nil {
		return nil, e
	}

	return rawData, nil
}
