package main

import (
	"bitbucket.org/nazwa/free-postcode-lottery-checker/config"
	"bitbucket.org/nazwa/free-postcode-lottery-checker/fpl"
	"fmt"
	"github.com/kardianos/service"
	"github.com/keighl/mandrill"
	"github.com/stvp/rollbar"
	"runtime"
	"strings"
	"time"
)

// Program structures.
// Define Start and Stop methods.
type program struct {
	exit chan struct{}

	client *fpl.FPLClient
	mailer *mandrill.Client

	daily    *fpl.Daily
	stockpot *fpl.Stockpot
	survey   *fpl.Survey
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) init() {

	config.LoadConfig("/config.json")

	// Start with external services
	p.mailer = mandrill.ClientWithKey(config.Config.Services.Mandrill.Key)
	rollbar.Environment = config.Config.Services.Rollbar.Environment
	rollbar.Token = config.Config.Services.Rollbar.Token

	// Now the games
	p.daily = &fpl.Daily{}
	p.stockpot = &fpl.Stockpot{}
	p.survey = &fpl.Survey{}

	// And finally the main worker
	p.client = fpl.NewClient(p.daily, p.stockpot, p.survey)
	if err := p.client.Login(); err != nil {
		logger.Error(err)
	}

}

func (p *program) run() error {
	// Allow the main thread to finish
	// This prevents the service from being terminated
	runtime.Gosched()

	p.init()

	ticker := time.Tick(3 * time.Hour)

	for {
		logger.Info("Fetching new results")
		if err := p.client.Run(); err != nil {
			logger.Error(err)
			rollbar.Error(rollbar.ERR, err)
		}
		if p.client.Changed() {
			logger.Info("New results detected. Sending email...")
			p.SendMail()
		}
		logger.Info("Going back to sleep")
		// Sleep until the next update time
		select {
		case <-ticker:
		}
	}

	return nil
}

func (p *program) SendMail() {

	win := p.client.CheckWin(config.Config.Target)
	templateVars := map[string]string{}

	message := &mandrill.Message{}
	if win {
		message.Subject = "WIN --- Postcode Lottery"
	} else {
		message.Subject = "Postcode Lottery"
	}
	message.InlineCSS = true
	message.Subaccount = "fpl"

	message.FromEmail = "chat@dimes.io"
	message.FromName = "FPL"

	message.AddRecipient("maciej@tidepayments.com", "Maciej", "to")

	if p.daily.Changed() {

		att := &mandrill.Attachment{
			Type:    "image/png",
			Name:    "daily",
			Content: p.daily.GetEncodedImage(),
		}
		templateVars["daily"] = "CHANGED"
		message.Attachments = make([]*mandrill.Attachment, 1)
		message.Attachments[0] = att
	}

	// Global vars

	if p.stockpot.Changed() {
		templateVars["stockpot"] = strings.Join(p.stockpot.GetPostcodes(), "<br>")
	}
	if p.survey.Changed() {
		templateVars["survey"] = p.survey.GetPostcode()
	}
	message.GlobalMergeVars = mandrill.MapToVars(templateVars)

	templateContent := map[string]string{}
	responses, err := p.mailer.MessagesSendTemplate(message, "fpl", templateContent)

	if err != nil {
		rollbar.Error(rollbar.ERR, err)
		if config.Config.Debug {
			logger.Error(err)
		}
	}

	for _, response := range responses {
		if strings.EqualFold(response.Status, "rejected") || strings.EqualFold(response.Status, "invalid") {
			apiErrorField := &rollbar.Field{Name: "Response", Data: response}
			rollbar.Error(rollbar.ERR, fmt.Errorf("Email send failed"), apiErrorField)
			if config.Config.Debug {
				logger.Error(response)
			}
		}

	}
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}
