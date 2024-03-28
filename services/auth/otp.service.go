package auth_service

import (
	"crypto/rand"
	"fmt"
	"net/smtp"
	"os"

	"github.com/ArifKobel/creator-tools/services/database"
	"github.com/ArifKobel/creator-tools/services/database/schemas"
)

func SendOtp(email string, otp string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	message := []byte("Your OTP is " + otp)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, message)
	if err != nil {
		fmt.Println(from)
		fmt.Println(password)
		fmt.Println(smtpHost)
		fmt.Println(smtpPort)
		return err
	}

	fmt.Println("Email Sent!")
	return nil
}

func SaveOtp(email string, otp string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	db.Model(schemas.User{}).Where("email = ?", email).Update("otp", otp)
	return nil
}

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
