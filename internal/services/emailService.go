package services

import (
	"fmt"
	smtp "net/smtp"
	config "notification/internal/config"
	models "notification/internal/models"
)

type EmailService interface {
	SendEmail(req models.EmailRequest) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendEmail(req models.EmailRequest) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.cfg.SMTPUser, req.To, req.Subject, req.Body,
	))

	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	return smtp.SendMail(addr, auth, s.cfg.SMTPUser, []string{req.To}, msg)
}
