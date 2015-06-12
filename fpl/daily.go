package fpl

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

type Daily struct {
	lastImage []byte
}

func (d *Daily) GetUrl() string {
	return "http://freepostcodelottery.com/speech/1.php?s=4"
}

func (d *Daily) Parse(response *http.Response) error {
	var err error
	d.lastImage, err = ioutil.ReadAll(response.Body)

	return err
}

func (d *Daily) GetEncodedImage() string {
	return base64.StdEncoding.EncodeToString(d.lastImage)
}
