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
	ReplyTo  string // 添加 ReplyTo 字段, 可选
}

// NewSmtpPoster 现在接收 replyTo 参数
func NewSmtpPoster(host string, protocol bool, port int, username string, password string, from string, replyTo string) *SmtpPoster {
	return &SmtpPoster{
		Host:     host,
		Protocol: protocol,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		ReplyTo:  replyTo,
	}
}

func (s *SmtpPoster) Valid() bool {
	// 保持原有的 Valid 逻辑，ReplyTo 不是必须的
	return s.Host != "" && s.Port > 0 && s.Port <= 65535 && s.Username != "" && s.Password != "" && s.From != ""
}

func (s *SmtpPoster) SendMail(to string, subject string, body string) error {
	if !s.Valid() {
		return fmt.Errorf("smtp not configured properly")
	}

	message := gomail.NewMessage()

	var from string
	if strings.Contains(s.Username, "@") {
		from = s.From
	} else {
		from = fmt.Sprintf("%s <%s>", s.Username, s.From)
	}
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// 设置 Reply-To, 优先使用配置的 ReplyTo，如果没有配置，则使用硬编码的默认值
	replyTo := s.ReplyTo
	if replyTo == "" {
		replyTo = "godisljr@163.com" // 替换成你的默认 Reply-To 地址
	}
	message.SetHeader("Reply-To", replyTo)
	print(replyTo)

	dialer := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	if s.Protocol {
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Host,
		}
	} else {
		dialer.SSL = true
	}

	switch {
	case strings.Contains(s.Host, "outlook"):
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "qq"):
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "office365"):
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "resend"):
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.Host,
		}
	case strings.Contains(s.Host, "tencent"):
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
	tmpl, err := template.New(filename).ParseFiles(fmt.Sprintf("utils/templates/%s", filename))
	if err != nil {
		return "", fmt.Errorf("template parsing error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

func (s *SmtpPoster) RenderMail(filename string, data interface{}, to string, subject string) error {
	body, err := s.RenderTemplate(filename, data)
	if err != nil {
		return fmt.Errorf("template rendering failed: %w", err)
	}

	return s.SendMail(to, subject, body)
}
