package mail

import (
	"fmt"
	"net/smtp"
	"purple/links/configs"

	"crypto/tls"

	"github.com/jordan-wright/email"
)

func SendVerifyEmail(e_mail string, text string, conf *configs.Config) {
	e := email.NewEmail()
	e.From = conf.SMTPConfig.Email
	e.To = []string{e_mail}
	e.Subject = "Подтверждение от сервиса Go Purple Advanced"
	e.Text = []byte(text)
	e.SendWithTLS(
		fmt.Sprintf("%s:%s", conf.SMTPConfig.Address, conf.SMTPConfig.Port),
		smtp.PlainAuth("", conf.SMTPConfig.Email, conf.SMTPConfig.Password, conf.SMTPConfig.Address),
		&tls.Config{
			ServerName: conf.SMTPConfig.Address,
		},
	)
}
