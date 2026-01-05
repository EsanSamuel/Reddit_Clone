package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/jobs/workers"
	"github.com/EsanSamuel/Reddit_Clone/models"
	"github.com/EsanSamuel/Reddit_Clone/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding payload request", "Details": err.Error()})
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			fmt.Println("Error hashing password")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user.UserId = bson.NewObjectID().Hex()
		user.Password = hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.EmailVerified = false

		userCount, err := database.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error counting user"})
			return
		}

		if userCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
			return
		}

		verificationToken, err := utils.GenerateVerificationOrResetToken()

		if err != nil {
			fmt.Println("Error generating verification token")
			return
		}

		user.VerficationToken = verificationToken

		_, err = utils.SendVerificationEmail(user.Email, user.VerficationToken)
		if err != nil {
			fmt.Println("Error sending verification mail")
		}

		result, err := database.UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": result})

	}
}

func VerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		token := c.Query("token")

		filter := bson.M{"verification_token": token}

		err := database.UserCollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error finding user", "details": err.Error()})
			return
		}

		updateData := bson.M{
			"$set": bson.M{
				"email_verified": true,
			}, "$unset": bson.M{"verification_token": ""},
		}

		result, err := database.UserCollection.UpdateOne(ctx, bson.M{"user_id": user.UserId}, updateData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error verfying user", "details": err.Error()})
			return
		}

		if result.Acknowledged {
			workers.SendEmailQueue(user.Email, user.UserId)
		}

		c.JSON(http.StatusOK, gin.H{"message": "User verified"})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userLogin models.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding request body", "details": err.Error()})
			return
		}

		var user models.User

		err := database.UserCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user"})
			return
		}

		if user.EmailVerified == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email is not verified"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
			return
		}

		token, refreshToken, err := utils.GenerateTokens(user.FirstName, user.LastName, user.Email, user.Role, user.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token", "details": err.Error()})
			return
		}

		err = utils.UpdateTokens(token, refreshToken, user.UserId, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error updating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "User": models.UserDTO{
			UserId:       user.UserId,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			Role:         user.Role,
			Token:        token,
			RefreshToken: refreshToken,
		}})

	}
}

func ResetPasswordRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var resetUser models.ForgetPasswordRequestDTO

		if err := c.ShouldBindJSON(&resetUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding forgot password payload", "details": err.Error()})
			return
		}

		resetToken, err := utils.GenerateVerificationOrResetToken()
		if err != nil {
			fmt.Println("Error generating reset token")
			return
		}

		var user models.User

		err = database.UserCollection.FindOne(ctx, bson.M{"email": resetUser.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user"})
			return
		}

		updateData := bson.M{
			"$set": bson.M{
				"reset_token": resetToken,
			},
		}

		_, err = database.UserCollection.UpdateOne(ctx, bson.M{"user_id": user.UserId}, updateData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating reset token", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reset token generated"})
	}
}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var reset models.ForgetPasswordDTO

		if err := c.ShouldBindJSON(&reset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding forgot password payload", "details": err.Error()})
			return
		}

		var user models.User

		filter := bson.M{"reset_token": token}

		err := database.UserCollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user"})
			return
		}

		hashedPassword, err := utils.HashPassword(reset.Password)

		updateData := bson.M{
			"$set": bson.M{
				"password": hashedPassword,
			}, "$unset": bson.M{
				"reset_token": "",
			},
		}

		_, err = database.UserCollection.UpdateOne(ctx, bson.M{"user_id": user.UserId}, updateData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user password", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated!"})
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var users []models.User

		filter := bson.M{}
		findOptions := options.Find()

		// Search User
		if s := strings.TrimSpace(c.Query("search")); s != "" {
			safe := regexp.QuoteMeta(s)

			filter = bson.M{
				"$or": []bson.M{
					{
						"first_name": bson.M{
							"$regex":   safe,
							"$options": "i",
						},
					},
					{
						"last_name": bson.M{
							"$regex":   safe,
							"$options": "i",
						},
					},
					{
						"role": bson.M{
							"$regex":   safe,
							"$options": "i",
						},
					},
				},
			}
		}

		// Sort User
		if sort := strings.TrimSpace(c.Query("sort")); sort != "" {
			switch sort {
			case "asc":
				findOptions.SetSort(bson.D{{"created_at", 1}})
			case "desc":
				findOptions.SetSort(bson.D{{"created_at", -1}})
			}

		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 9

		findOptions.SetSkip((int64(page) - 1) * int64(perPage))
		findOptions.SetLimit(int64(perPage))

		cursor, err := database.UserCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching users", "details": err.Error()})
			return
		}

		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding users", "details": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("userId")

		var user models.User

		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is not provided"})
			return
		}

		filter := bson.M{"user_id": userId}

		err := database.UserCollection.FindOne(ctx, filter).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching user", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

/*func UploadFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		form, err := c.MultipartForm()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error getting files", "details": err.Error()})
			return
		}

		files := form.File["files"]

		urls := utils.UploadFiles(c, files)

		c.JSON(http.StatusCreated, urls)
	}
}*/

func UploadAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user_id := c.Param("userId")

		file, err := c.FormFile("avatar")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error uploading file", "details": err.Error()})
			return
		}

		src, err := file.Open()
		fileByte, err := io.ReadAll(src)

		if !utils.IsFileImage(fileByte) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is not a image", "details": err.Error()})
			return
		}

		AvatarUrl, err := utils.UploadSingleFileToS3(src, file.Filename)

		updateAvatar := bson.M{
			"$set": bson.M{
				"avatar": AvatarUrl,
			},
		}

		_, err = database.UserCollection.UpdateOne(ctx, bson.M{"user_id": user_id}, updateAvatar)

		c.JSON(http.StatusOK, AvatarUrl)

	}
}
