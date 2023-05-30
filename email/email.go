package email

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Recipient struct {
	FirstName string
	LastName  string
	Email     string
}

func SendEmail() {
	from := mail.NewEmail("Example User", "John")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "example@email")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
}

func SendWelcomeEmail(recipient Recipient) {
	from := mail.NewEmail("Coeus Education", "noreply@coeuseducation.com")
	subject := "Welcome to Coeus Education"
	to := mail.NewEmail(recipient.FirstName, recipient.Email)
	plainTextContent := fmt.Sprintf("Hello %s, welcome to Coeus Education!", recipient.FirstName)
	htmlContent := fmt.Sprintf("<strong>Hello %s, welcome to Coeus Education!</strong>", recipient.FirstName)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	sendEmail(message)
}

func SendForgotPasswordEmail(userEmail, name string) string {
	organizationEmail := os.Getenv("SENDGRID_ORGANIZATION_EMAIL")

	code := generateRandomCode()
	from := mail.NewEmail("Coeus Education", organizationEmail)
	subject := "Coeus Education - Password Reset"
	to := mail.NewEmail(name, userEmail)

	plainTextContent := fmt.Sprintf("Your password reset code is: %s", code)
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<style>
  body {
    font-family: Arial, sans-serif;
  }
  .container {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    background-color: #f9f9f9;
    border-radius: 10px;
  }
  .header {
    text-align: center;
  }
  .content {
    margin-top: 30px;
    text-align: center;
  }
  .reset-code {
    font-size: 24px;
    font-weight: bold;
    color: #3c3c3c;
    margin-top: 20px;
  }
  .code-box {
    display: inline-block;
    padding: 10px 15px;
    background-color: #eee;
    border: 1px solid #ccc;
    border-radius: 5px;
  }
  .footer {
    margin-top: 40px;
    text-align: center;
    font-size: 14px;
    color: #777;
  }
</style>
</head>
<body>
<div class="container">
  <div class="header">
    <img src="https://coeus.education/images/coeus-banner.png" alt="Coeus Education" style="width: 80% !important;">
  </div>
  <div class="content">
    <p>Dear user,</p>
    <p>We have received a request to reset your password for your Coeus Education account. Please use the code below to reset your password:</p>
    <p class="reset-code">Copy this code: <span class="code-box">%[1]s</span></p>
    <p>Please select and copy the code above, then paste it into the appropriate field on the password reset page. If you didn't request this password reset, please ignore this email. The code will expire in 24 hours.</p>
  </div>
  <div class="footer">
    <p>Best regards,</p>
    <p>The Coeus Education Team</p>
  </div>
</div>
</body>
</html>
`, code)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	sendEmail(message)

	return code
}

func sendEmail(message *mail.SGMailV3) {

	apiKey := os.Getenv("SENDGRID_API_KEY")
	client := sendgrid.NewSendClient(apiKey)
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
	// Useful for debugging
	// else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
}

func generateRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	pin := rand.Intn(1000000)
	return fmt.Sprintf("%06d", pin)
}
