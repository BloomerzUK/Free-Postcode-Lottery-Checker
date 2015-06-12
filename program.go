package main

import (
	"bitbucket.org/nazwa/free-postcode-lottery-checker/fpl"
	"github.com/kardianos/service"
	"github.com/keighl/mandrill"
	"log"
	"runtime"
	"strings"
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

	p.daily = &fpl.Daily{}
	p.stockpot = &fpl.Stockpot{}
	p.survey = &fpl.Survey{}

	p.client = fpl.NewClient(p.daily, p.stockpot, p.survey)
	p.client.Login()

	p.mailer = mandrill.ClientWithKey("o_3WECMMBJ3qouJqMKi0Fg")
}

func (p *program) run() error {
	// Allow the main thread to finish
	// This prevents the service from being terminated
	runtime.Gosched()

	p.init()

	p.client.Run()
	p.SendMail()

	return nil
}

func (p *program) SendMail() {
	att := &mandrill.Attachment{
		Type:    "image/png",
		Name:    "daily",
		Content: p.daily.GetEncodedImage(),
	}

	message := &mandrill.Message{}
	message.Subject = "Postcode Lottery"
	message.InlineCSS = true
	message.Subaccount = "fpl"

	message.Attachments = make([]*mandrill.Attachment, 1)
	message.Attachments[0] = att

	message.AddRecipient("zorleq@hotmail.com", "Maciej", "to")

	// Global vars

	templateVars := map[string]string{}
	templateVars["stockpot"] = strings.Join(p.stockpot.GetPostcodes(), "<br>")
	templateVars["survey"] = strings.Join(p.survey.GetPostcodes(), "<br>")
	message.GlobalMergeVars = mandrill.MapToVars(templateVars)

	templateContent := map[string]string{}
	responses, err := p.mailer.MessagesSendTemplate(message, "fpl", templateContent)

	if err != nil {
		log.Println(err)
	}

	for _, response := range responses {
		if strings.EqualFold(response.Status, "rejected") || strings.EqualFold(response.Status, "invalid") {
			log.Println(response.RejectionReason)
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}
