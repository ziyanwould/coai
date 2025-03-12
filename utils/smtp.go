package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"strings"
	"text/template"

	"gopkg.in/gomail.v2"
)

type SmtpPoster struct {
	Host     string
	Protocol bool
	Port     int
	Username string
	Password string
	From     string
	//ReplyTo  string // 先注释掉，后面再加
}

func NewSmtpPoster(host string, protocol bool, port int, username string, password string, from string) *SmtpPoster {
	return &SmtpPoster{
		Host:     host,
		Protocol: protocol,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		//ReplyTo:  replyTo, // 注释掉
	}
}

func (s *SmtpPoster) Valid() bool {
	return s.Host != "" && s.Port > 0 && s.Port <= 65535 && s.Username != "" && s.Password != "" && s.From != ""
}

func (s *SmtpPoster) SendMail(to string, subject string, body string) error {
	if !s.Valid() {
		return fmt.Errorf("smtp not configured properly")
	}

	// Create gomail message object
	message := gomail.NewMessage()

	// Determine sender address
	var from string
	if strings.Contains(s.Username, "@") {
		from = s.From
	} else {
		from = fmt.Sprintf("%s <%s>", s.Username, s.From)
	}
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body) // 统一使用 text/html

	// --- 硬编码 Reply-To ---
	replyTo := "godisljr@163.com" // 替换成你的测试邮箱
	if replyTo != "" {            // 即使硬编码，也建议保留这个检查
		message.SetHeader("Reply-To", replyTo)
	}

	dialer := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	// TLS/SSL configuration
	if s.Protocol {
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Host,
		}
	} else {
		dialer.SSL = true
	}

	// Provider-specific handling
	switch {
	case strings.Contains(s.Host, "outlook"):
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "qq"):
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // QQ邮箱需要
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "office365"):
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "resend"): // resend 好像需要 InsecureSkipVerify: true
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "tencent"): // 腾讯企业邮箱
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.Host,
		}
	}

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("sent mail failed: %s", err.Error())
	}

	return nil
}

func (s *SmtpPoster) RenderTemplate(filename string, data interface{}) (string, error) {
	tmpl, err := template.New(filename).ParseFiles(fmt.Sprintf("utils/templates/%s", filename)) //更正了templates路径
	if err != nil {
		return "", fmt.Errorf("template parsing error: %w", err) // 更好的错误信息
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil { // 更简单，不需要 ExecuteTemplate
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

func (s *SmtpPoster) RenderMail(filename string, data interface{}, to string, subject string) error {
	body, err := s.RenderTemplate(filename, data)
	if err != nil {
		return fmt.Errorf("template rendering failed: %w", err) // 更清晰的错误
	}

	return s.SendMail(to, subject, body)
}
