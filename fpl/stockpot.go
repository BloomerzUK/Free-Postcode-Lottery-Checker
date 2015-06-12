package fpl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Stockpot struct {
	postcodes []string
}

func (d *Stockpot) GetUrl() string {
	return "http://freepostcodelottery.com/stackpot/"
}

func (d *Stockpot) Parse(response *http.Response) error {
	pattern := regexp.MustCompile("<span>(.*)</span>")

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

func (d *Stockpot) GetPostcodes() []string {
	return d.postcodes
}
