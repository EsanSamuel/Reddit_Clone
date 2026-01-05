package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/resend/resend-go/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	Firstname string
	LastName  string
	Email     string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
var JWT_SECRET_REFRESH_KEY = os.Getenv("JWT_SECRET_REFRESH_KEY")

func GenerateTokens(firstname string, lastname string, email string, role string, user_id string) (string, string, error) {
	claims := &SignedDetails{
		Firstname: firstname,
		LastName:  lastname,
		Email:     email,
		Role:      role,
		UserId:    user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Reddit_Clone",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(JWT_SECRET_KEY))

	if err != nil {
		fmt.Println("Error generating jwt token")
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		Firstname: firstname,
		LastName:  lastname,
		Email:     email,
		Role:      role,
		UserId:    user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Reddit_Clone",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(JWT_SECRET_REFRESH_KEY))

	if err != nil {
		fmt.Println("Error generating jwt token")
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil

}

func UpdateTokens(token string, refreshToken string, userId string, c *gin.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updateAt,
		},
	}

	_, err := database.UserCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error updating token"})
		return err
	}

	return nil

}

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password")
		return "", err
	}

	return string(hashPassword), nil
}

func GenerateVerificationOrResetToken() (string, error) {
	b := make([]byte, 35)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error generating token")
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func SendVerificationEmail(email string, verificationToken string) (string, error) {
	url := "http://localhost:3000/verify-email?token=" + verificationToken
	RESEND_API_KEY := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(RESEND_API_KEY)

	params := &resend.SendEmailRequest{
		From: "Acme <noreply@mikaelsoninitiative.org>",
		To:   []string{email},
		Html: `<div style="max-width: 500px; margin: 0 auto; font-family: Arial, sans-serif; background-color: #ffffff; padding: 30px; border-radius: 8px; border: 1px solid #e5e7eb;">
  
  <h2 style="color: #111827; text-align: center; margin-bottom: 10px;">
    Confirm Your Signup
  </h2>

  <p style="color: #374151; font-size: 15px; text-align: center;">
    Hey there 👋
  </p>

  <p style="color: #374151; font-size: 15px; text-align: center; line-height: 1.5;">
    Thanks for joining <b>Reddit</b>! Please confirm your email address to activate your account.
  </p>

  <div style="text-align: center; margin: 30px 0;">
    <a href="` + url + `"
       style="
         background-color: #2563eb;
         color: #ffffff;
         padding: 14px 30px;
         text-decoration: none;
         border-radius: 6px;
         font-weight: bold;
         display: inline-block;
         font-size: 16px;
       ">
      Confirm Email
    </a>
  </div>

  <p style="color: #6b7280; font-size: 14px; text-align: center; line-height: 1.4;">
    If you didn’t sign up, you can safely ignore this email.
  </p>

</div>`,
		Subject: "Hello from Golang",
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(sent.Id)
	return sent.Id, nil
}

func GetAuthToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		fmt.Println("There is no authorization token")
		return "", fmt.Errorf("Authorization token not found")
	}

	authToken := strings.Split(authHeader, " ")[1]
	fmt.Println("TokenString:", authToken)

	return authToken, nil
}

func VerifyAuthToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired!")
	}

	return claims, nil
}

func uploadToS3(file multipart.File, filename string, resultChan chan<- string) (string, error) {
	defer file.Close()

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	UPLOAD_URL := os.Getenv("UPLOAD_URL")
	AWS_ENDPOINT := os.Getenv("AWS_ENDPOINT")
	AWS_REGION := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(AWS_REGION),
		Endpoint: aws.String(AWS_ENDPOINT),
		Credentials: credentials.NewStaticCredentials(
			accessKey, secretKey, "",
		),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	s3Client := s3.New(sess)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	})

	url := UPLOAD_URL + filename

	resultChan <- url

	return url, err
}

func UploadFiles(c *gin.Context, files []*multipart.FileHeader) []string {
	result := make(chan string, len(files))

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error opening file", "details": err.Error()})
			continue
		}
		go uploadToS3(src, file.Filename, result)
	}

	var urls []string
	for i := 0; i < len(files); i++ {
		url := <-result
		urls = append(urls, url)
	}

	fmt.Println("All uploaded urls:", urls)
	return urls
}

func UploadSingleFileToS3(file multipart.File, filename string) (string, error) {
	defer file.Close()

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	UPLOAD_URL := os.Getenv("UPLOAD_URL")
	AWS_ENDPOINT := os.Getenv("AWS_ENDPOINT")
	AWS_REGION := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(AWS_REGION),
		Endpoint: aws.String(AWS_ENDPOINT),
		Credentials: credentials.NewStaticCredentials(
			accessKey, secretKey, "",
		),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	s3Client := s3.New(sess)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	})

	url := UPLOAD_URL + filename

	return url, err
}

func IsFileImage(file []byte) bool {
	fileType := http.DetectContentType(file)
	return strings.HasPrefix(fileType, "image/")
}
