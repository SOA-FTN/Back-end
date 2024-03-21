package service

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/smtp"
	"stakeholders/model"
	"stakeholders/repo"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	AuthRepo *repo.AuthRepository
}

var jwtKey = string ("explorer_secret_key")

//Authentifikacija
func (service *AuthService) Authentication (credentials *model.Credentials) (*model.User,error) {
	user,err := service.AuthRepo.Authentication(credentials)
	return user,err
}
//Generisanje tokena
func (service *AuthService) GenerateToken(user *model.User) (string, error) {
    expirationTime := time.Now().Add(time.Minute * 60 * 24) 

    claims := jwt.MapClaims{
        "id":       user.ID,
        "username": user.UserName,
        "role":     user.GetRoleName(),
        "exp":      expirationTime.Unix(),
        "iss":      "explorer",
        "aud":      "explorer-front.com",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jwtKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func(service *AuthService) SendVerificationMail(registration *model.Registration, verificationToken string) (error) {
    // Sender's email configuration
    from := "pswexplorer@gmail.com"
    password := "eqvw sehe rgzf nzzp"
    smtpHost := "smtp.gmail.com"
    smtpPort := "587" // This is typically the port for SMTP submission

    // Recipient email address
    to := registration.Email

    // Message to be sent
    message := []byte("Subject: Email confirmation!\r\n" +
        "From: " + from + "\r\n" +
        "To: " + to + "\r\n" +
        "\r\n" +
        "Click the following link to verify your email: \n http://localhost:4200/confirm?token=" + verificationToken + "\r\n")

    // Authentication
    auth := smtp.PlainAuth("", from, password, smtpHost)

    // Sending email
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
    if err != nil {
        log.Fatal(err)
        return err
    }

    //log.Println("Email sent successfully!")
    return nil
}


func(service *AuthService) GenerateUniqueVerificationToken() string {
    tokenBytes := make([]byte, 32)
    _, err := rand.Read(tokenBytes)
    if err != nil {
        log.Fatal(err)
    }

    verificationToken := hex.EncodeToString(tokenBytes)

    return verificationToken
}
