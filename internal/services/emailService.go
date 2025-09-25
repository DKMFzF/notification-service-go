package services

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"notification/internal/config"
	"notification/internal/models"
	"notification/pkg/services"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) services.Notifier[models.EmailRequest] {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) Send(req models.EmailRequest) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.cfg.SMTPUser, req.To, req.Subject, req.Body,
	))
	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	return smtp.SendMail(addr, auth, s.cfg.SMTPUser, []string{req.To}, msg)
}

func EmailConverter(data []byte) (models.EmailRequest, error) {
	var req models.EmailRequest
	err := json.Unmarshal(data, &req)
	return req, err
}

// register in storage services
func init() {
	Register("email", Factory{
		NewService: func(cfg *config.Config) services.Notifier[any] {
			return NotifierAdapter[models.EmailRequest]{NewEmailService(cfg)}
		},
		Converter: func(b []byte) (any, error) {
			return EmailConverter(b)
		},
	})
}

type NotifierAdapter[T any] struct{ n services.Notifier[T] }

func (a NotifierAdapter[T]) Send(v any) error {
	req, ok := v.(T)
	if !ok {
		return fmt.Errorf("wrong type for notifier")
	}
	return a.n.Send(req)
}
