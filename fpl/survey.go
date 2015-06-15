package fpl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Survey struct {
	postcode string
	oldCode  string
	temp     string
}

func (d *Survey) GetUrl() string {
	return "http://freepostcodelottery.com/survey-draw/"
}

func (d *Survey) Parse(response *http.Response) error {
	pattern := regexp.MustCompile("<p style=\"width:100%;margin:18px 0 24px;font-size:56px;font-weight:bold;text-align:center\">(.*)</p>")

	html, _ := ioutil.ReadAll(response.Body)

	d.oldCode = d.postcode
	postcode := pattern.FindSubmatch(html)
	if postcode == nil {
		return fmt.Errorf("No matches found")
	}

	d.postcode = string(postcode[1])

	return nil
}

func (d *Survey) GetPostcode() string {
	return d.postcode
}

func (d *Survey) Changed() bool {
	return !strings.EqualFold(d.postcode, d.oldCode)
}

func (d *Survey) Check(postcode string) bool {
	return strings.EqualFold(postcode, d.postcode)
}
