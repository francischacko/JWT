package controller

import (
	"gin/initializers"
	"gin/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Signup(C *gin.Context) {
	// get the email and password required
	var body struct {
		Email    string
		Password string
	}

	if C.Bind(&body) != nil {
		C.JSON(400, gin.H{
			"error": "failed to load",
		})
		return
	}
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	// create user
	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		C.JSON(400, gin.H{
			"error": "failed to create user",
		})
		return
	}
	//respond
	C.JSON(http.StatusOK, gin.H{
		"message": "user registered",
	})
}

func Login(C *gin.Context) {
	// get email and password required off the body
	var body struct {
		Email    string
		Password string
	}

	if C.Bind(&body) != nil {
		C.JSON(400, gin.H{
			"error": "failed to load",
		})
		return
	}

	// lookup requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		C.JSON(400, gin.H{
			"error": "invalid user or password",
		})
		return
	}

	// compare sent in hash  password with saved user hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		C.JSON(400, gin.H{
			"error": "invalid user or password",
		})
		return
	}
	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed to create token",
		})
		return
	}
	// send it back
	C.SetSameSite(http.SameSiteLaxMode)
	C.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	C.JSON(http.StatusOK, gin.H{
		"message": user,
		"token":   tokenString,
	})

	// return
}

func Validate(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "user is logged in.",
	})

}
