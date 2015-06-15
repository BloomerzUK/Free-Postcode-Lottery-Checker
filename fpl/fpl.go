package fpl

import (
	"net/http"
	"net/url"
	"time"
)

type FPLInterface interface {
	GetUrl() string
	Changed() bool
	Check(string) bool
	Parse(*http.Response) error
}

type FPLClient struct {
	client  *http.Client
	request *http.Request
	games   []FPLInterface
	time    time.Time
}

func NewClient(games ...FPLInterface) *FPLClient {
	fpl := &FPLClient{}
	fpl.client = &http.Client{}

	fpl.request, _ = http.NewRequest("GET", "", nil)

	fpl.request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	fpl.request.Header.Add("Cache-Control", "no-cache")
	fpl.request.Header.Add("Accept-Encoding", "gzip, deflate, sdch")
	fpl.request.Header.Add("Accept-Language", "en-GB,en;q=0.8,pl;q=0.6,es;q=0.4")
	fpl.request.Header.Add("Cache-Control", "no-cache")
	fpl.request.Header.Add("Connection", "keep-alive")
	fpl.request.Header.Add("DNT", "1")
	fpl.request.Header.Add("Host", "freepostcodelottery.com")
	fpl.request.Header.Add("Pragma", "no-cache")
	fpl.request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36")

	fpl.time = time.Now().AddDate(0, 0, -1)

	fpl.games = games

	return fpl
}

func (f *FPLClient) Fetch(targetUrl string) (*http.Response, error) {
	f.request.URL, _ = url.Parse(targetUrl)

	diff := time.Now().Sub(f.time)
	if diff < 5*time.Second {
		time.Sleep(5*time.Second - diff)
	}
	defer func() {
		f.time = time.Now()
	}()

	return f.client.Do(f.request)
}

func (f *FPLClient) Login() error {

	resp, err := f.Fetch("http://freepostcodelottery.com/?reminder=b4ff6706-108b-11e5-ad49-00163ee58471")
	if err != nil {
		return err
	}

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		f.request.AddCookie(cookie)
	}

	return nil
}

func (f *FPLClient) Run() error {
	for _, game := range f.games {
		if resp, err := f.Fetch(game.GetUrl()); err == nil {
			if err = game.Parse(resp); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (f *FPLClient) Changed() bool {
	for _, game := range f.games {
		res := game.Changed()
		if res {
			return true
		}
	}
	return false
}

func (f *FPLClient) CheckWin(postcode string) bool {

	for _, game := range f.games {
		res := game.Check(postcode)
		if res {
			return true
		}
	}
	return false
}
