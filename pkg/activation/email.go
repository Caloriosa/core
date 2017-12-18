package activation

import (
	"bytes"
	"core/types"
	"crypto/tls"
	"github.com/golang/glog"
	"html/template"
	"gopkg.in/gomail.v2"
	"core/pkg/config"
)

func SendValidationEmail(user *types.User) error {
	t, err := template.New("validation").ParseFiles("static/email_validation.html")
	if err != nil {
		glog.Error("Error creating a validation email: ", err)
		return err
	}

	buffer := new(bytes.Buffer)

	err = t.ExecuteTemplate(buffer, "email_validation.html", user)
	if err != nil {
		glog.Error("Error executing a new validation email: ", err)
	}

	//err = InsecureSendMail("10.0.0.90:925", "caloriosa@mail.foxiehost.lan", "caloriosa@victorianfox.com", buffer.Bytes(),
	//	smtp.PlainAuth("", "caloriosa", "caloriosa", "10.0.0.90"))

	dialer := gomail.NewDialer(config.LoadedConfig.Email.SmtpHost, config.LoadedConfig.Email.SmtpPort, config.LoadedConfig.Email.SmtpUser, config.LoadedConfig.Email.SmtpPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify:true, ServerName:config.LoadedConfig.Email.SmtpHost}
	email := gomail.NewMessage()
	email.SetHeader("From", config.LoadedConfig.Email.EmailFrom)
	email.SetHeader("To", config.LoadedConfig.Dev.TestEmailTo)
	email.SetHeader("Subject", "Activation at Caloriosa")
	email.SetBody("text/html", buffer.String())

	if err = dialer.DialAndSend(email); err != nil {
		glog.Error("Error sending email: ", err)
		return err
	} else {
		glog.Info("Activation email sent")
	}

	return nil
}