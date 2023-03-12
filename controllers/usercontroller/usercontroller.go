package usercontroller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/yurdigrfnn/api-todo-auth/initializers"
	"github.com/yurdigrfnn/api-todo-auth/models"
	"golang.org/x/crypto/bcrypt"
)

// response json
type SignupResponse struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

type SigninResponse struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// controller function
func Signup(c *gin.Context) {

	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, SignupResponse{
			IsError: true,
			Message: "failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, SignupResponse{
			IsError: true,
			Message: "failed to hash",
		})

		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, SignupResponse{
			IsError: true,
			Message: "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, SignupResponse{
		IsError: false,
		Message: "succes to create user",
	})
}

func Signin(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, SigninResponse{
			IsError: true,
			Message: "failed to read body",
		})
		return
	}

	// req user from http body
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, SigninResponse{
			IsError: true,
			Message: "invalid email & password",
		})
		return
	}

	// compare pass

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, SigninResponse{
			IsError: true,
			Message: "error email & password",
		})
		return
	}

	//generate jwt

	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusBadRequest, SigninResponse{
			IsError: true,
			Message: "create token invalid",
		})
		return
	}
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Path:     "/",
		Domain:   "http://localhost:8000",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, SigninResponse{
		IsError: false,
		Message: "Successful login",
		Token:   tokenString,
	})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"isError": false,
		"data":    user,
	})
}

func Logout(c *gin.Context) {
	pastTime := time.Now().Add(-time.Hour) // set the expiration time in the past
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  pastTime,
		Path:     "/",
		Domain:   "http://localhost:8000",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(http.StatusOK, gin.H{
		"isError": false,
		"message": "Logout successful",
	})
}
