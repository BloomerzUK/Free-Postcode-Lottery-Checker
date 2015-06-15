package fpl

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

type Daily struct {
	oldImage []byte
	newImage []byte
}

func (d *Daily) GetUrl() string {
	return "http://freepostcodelottery.com/speech/1.php?s=4"
}

func (d *Daily) Parse(response *http.Response) error {
	var err error
	d.oldImage = d.newImage
	d.newImage, err = ioutil.ReadAll(response.Body)

	return err
}

func (d *Daily) GetEncodedImage() string {
	return base64.StdEncoding.EncodeToString(d.newImage)
}

func (d *Daily) Changed() bool {
	return !bytes.Equal(d.oldImage, d.newImage)
}

// We cant check the daily one, so just return failure
func (d *Daily) Check(postcide string) bool {
	return false
}
