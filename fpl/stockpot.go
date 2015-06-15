package fpl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Stockpot struct {
	oldCodes  []string
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

	// Save the old data
	d.oldCodes = make([]string, len(postcodes))
	copy(d.oldCodes, d.postcodes)

	// before making space for the new one
	d.postcodes = make([]string, len(postcodes))
	for _A, tag := range postcodes {
		d.postcodes[_A] = string(tag[1])
	}

	return nil
}

func (d *Stockpot) GetPostcodes() []string {
	return d.postcodes
}

func (d *Stockpot) Changed() bool {
	// Loops through the 'new' codes returns true if code is not in the old set
	for _, code := range d.postcodes {
		var found bool = false
		for _, old := range d.oldCodes {
			if old == code {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}
	return false
}

// We cant check the daily one, so just return failure
func (d *Stockpot) Check(postcode string) bool {
	for _, code := range d.postcodes {
		if strings.EqualFold(postcode, code) {
			return true
		}
	}
	return false
}
