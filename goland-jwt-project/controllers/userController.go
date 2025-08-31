package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Somu/golang-jwt-project/database"
	"github.com/Somu/golang-jwt-project/helpers"
	"github.com/Somu/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(hashedPassword, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		// err can be bcrypt.ErrMismatchedHashAndPassword or another error
		return false, err
	}
	return true, nil
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set default values for required fields if they're nil
		if user.FirstName == nil {
			defaultFirstName := "User"
			user.FirstName = &defaultFirstName
		}
		if user.LastName == nil {
			defaultLastName := "Default"
			user.LastName = &defaultLastName
		}
		if user.UserType == nil {
			defaultUserType := "USER"
			user.UserType = &defaultUserType
		}
		if user.Phone == nil {
			defaultPhone := "0000000000"
			user.Phone = &defaultPhone
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if email exists before proceeding
		if user.Email == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": *user.Email})
		if err != nil {
			log.Println("email count error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence while checking email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists with this email"})
			return
		}

		// Hash the password
		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": *user.Phone})
		if err != nil {
			log.Println("phone count error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence while checking phone number"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists with this phone number"})
			return
		}

		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.ID = primitive.NewObjectID()
		userID := user.ID.Hex()
		user.User_id = &userID
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.Email == nil || user.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
				return
			}
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching user details"})
			return
		}

		passwordIsValid, err := VerifyPassword(*foundUser.Password, *user.Password)
		if !passwordIsValid {
			log.Printf("Password verification failed for email %s: %v", *user.Email, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		token, refreshToken, err := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, *foundUser.User_id)
		if err != nil {
			log.Printf("Error generating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
			return
		}
		if err := helpers.UpdateAllTokens(token, refreshToken, *foundUser.User_id); err != nil {
			log.Printf("Failed to update tokens for user %s: %v", *foundUser.User_id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tokens"})
			return
		}

		// Send a curated response instead of the full user object to avoid leaking sensitive data like the password hash.
		c.JSON(http.StatusOK, gin.H{
			"user_id":       foundUser.User_id,
			"first_name":    foundUser.FirstName,
			"last_name":     foundUser.LastName,
			"email":         foundUser.Email,
			"user_type":     foundUser.UserType,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helpers.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recordPerPage"})
			return
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}

		startIndex := (page - 1) * recordPerPage

		// Fixed MongoDB aggregation pipeline
		skipStage := bson.D{{Key: "$skip", Value: startIndex}}
		limitStage := bson.D{{Key: "$limit", Value: recordPerPage}}
		countStage := bson.D{{Key: "$count", Value: "totalRecords"}}

		// Get total count
		countResult, err := userCollection.Aggregate(ctx, mongo.Pipeline{countStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting users"})
			return
		}

		var countData []bson.M
		if err = countResult.All(ctx, &countData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing count"})
			return
		}

		totalRecords := 0
		if len(countData) > 0 {
			if count, ok := countData[0]["totalRecords"].(int32); ok {
				totalRecords = int(count)
			}
		}

		// Get paginated users
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{skipStage, limitStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"totalRecords":  totalRecords,
			"users":         allusers,
			"page":          page,
			"recordPerPage": recordPerPage,
		})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.MatchUserType(c, userId); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
