package fpl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Survey struct {
	postcodes []string
	temp      string
}

func (d *Survey) GetUrl() string {
	return "http://freepostcodelottery.com/survey-draw/"
}

func (d *Survey) Parse(response *http.Response) error {
	pattern := regexp.MustCompile("<p style=\"width:100%;margin:18px 0 24px;font-size:56px;font-weight:bold;text-align:center\">(.*)</p>")

	html, _ := ioutil.ReadAll(response.Body)

	postcodes := pattern.FindAllSubmatch(html, 10)
	if postcodes == nil {
		return fmt.Errorf("No matches found")
	}
	d.postcodes = make([]string, len(postcodes))

	for _A, tag := range postcodes {
		d.postcodes[_A] = string(tag[1])
	}

	return nil
}

func (d *Survey) GetPostcodes() []string {
	return d.postcodes
}
