package main

import(
    "gopkg.in/gomail.v2"
    "crypto/tls"
    "log"
)
func Send(config *MailConfig, message *gomail.Message) {
    if config.Single {
        for _, o := range config.To {
            tos := []interface{}{}
            tos = append(tos, o)
            config2 := MailConfig{
                Provider: config.Provider,
                Host: config.Host,
                Port: config.Port,
                From: config.From,
                To: tos,
                Cc: config.Cc,
                Bcc: config.Bcc,
                Subject: config.Subject,
                Body: config.Body,
                Attachments: config.Attachments,
                Username: config.Username,
                Password: config.Password,
                PlainText: config.PlainText,
                Single: false,
                Thread: config.Thread,
                SkipAttachmentCheck: true,
            }
            if (config2.Thread) {
                go SendMail(&config2)
            } else {
                SendMail(&config2)
            }
        }
    } else {
        dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
        dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
        if err := dialer.DialAndSend(message); err != nil {
            log.Println("SEND_FAILED", err)
        }
    }
}

func SendMail(config *MailConfig) {
    Normalize(config)
    m := gomail.NewMessage()
    username := config.Username
    if config.From != "" {
        m.SetHeader("From", config.From)
        if username == "" {
            username = config.From
        }
    } else {
        config.From = username
    }

    config.Username = username

    headers := map[string][]string{}
    to := []string{}

    if len(config.To) > 0 {
        for _, o := range config.To {
            addr, _ := o.(map[string]interface{})
            a := addr["address"].(string)
            n := addr["name"].(string)
            if n == "" {
                to = append(to, a)
            } else {
                to = append(to, m.FormatAddress(a, n))
            }
        }
    }
    headers["To"] = to
    m.SetHeaders(headers)
    m.SetHeader("Subject", config.Subject)
    if config.PlainText {
        m.SetBody("text/plain", config.Body)
    } else {
        m.SetBody("text/html", config.Body)
    }

    if len(config.Attachments) > 0 {
        for _, o := range config.Attachments {
            fullpath, _ := o.(string)
            m.Attach(fullpath)
        }
    }

    if (config.Thread) {
        go Send(config, m)
    } else {
        Send(config, m)
    }
}
