package helper

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

func GetEmailConfig() EmailConfig {
	port, _ := strconv.Atoi(GetEnv("MAIL_PORT"))
	return EmailConfig{
		Host:     GetEnv("MAIL_HOST"),
		Port:     port,
		Username: GetEnv("MAIL_USERNAME"),
		Password: GetEnv("MAIL_PASSWORD"),
		From:     GetEnv("MAIL_FROM_ADDRESS"),
		FromName: GetEnv("MAIL_FROM_NAME"),
	}
}

func SendEmail(to []string, subject string, body string) error {
	config := GetEmailConfig()

	// Setup headers
	from := config.From
	if config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", config.FromName, config.From)
	}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Setup authentication
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// For SSL/TLS (port 465)
	if config.Port == 465 {
		tlsconfig := &tls.Config{
			ServerName: config.Host,
		}

		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("failed to dial: %v", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, config.Host)
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %v", err)
		}

		if err = client.Mail(config.From); err != nil {
			return fmt.Errorf("failed to set sender: %v", err)
		}

		for _, addr := range to {
			if err = client.Rcpt(addr); err != nil {
				return fmt.Errorf("failed to set recipient: %v", err)
			}
		}

		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to create data writer: %v", err)
		}

		_, err = writer.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write message: %v", err)
		}

		err = writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close writer: %v", err)
		}

		return client.Quit()
	}

	// For STARTTLS (port 587)
	return smtp.SendMail(addr, auth, config.From, to, []byte(message))
}

func SendPermitExpiryNotification(to []string, permitName string, permitNo string, expiryDate string, daysLeft int) error {
	subject := fmt.Sprintf("Reminder: Permit %s akan segera expired", permitName)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #dc2626; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9f9f9; padding: 20px; border: 1px solid #ddd; }
        .info-box { background-color: #fff; padding: 15px; margin: 15px 0; border-left: 4px solid #dc2626; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
        .warning { color: #dc2626; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Permit Expiry Notification</h1>
        </div>
        <div class="content">
            <p>Dear User,</p>
            <p>Ini adalah pengingat bahwa permit berikut akan segera expired:</p>
            
            <div class="info-box">
                <h3>Detail Permit:</h3>
                <p><strong>Nama Permit:</strong> %s</p>
                <p><strong>Nomor Permit:</strong> %s</p>
                <p><strong>Tanggal Expired:</strong> %s</p>
                <p class="warning">Sisa Waktu: %d hari</p>
            </div>
            
            <p>Harap segera melakukan perpanjangan atau pembaruan permit sebelum tanggal expired.</p>
            <p>Silakan login ke aplikasi Permit Management untuk informasi lebih lanjut.</p>
        </div>
        <div class="footer">
            <p>Email ini dikirim secara otomatis oleh Permit Management System.</p>
            <p>Mohon tidak membalas email ini.</p>
        </div>
    </div>
</body>
</html>
`, permitName, permitNo, expiryDate, daysLeft)

	return SendEmail(to, subject, body)
}
