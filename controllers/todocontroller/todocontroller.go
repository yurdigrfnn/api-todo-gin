package todocontroller

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/yurdigrfnn/api-todo-auth/initializers"
	"github.com/yurdigrfnn/api-todo-auth/models"
	"gorm.io/gorm"
)

type TodoIndexResponse struct {
	IsError    bool          `json:"isError"`
	Todos      []models.Todo `json:"todos"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"totalPages"`
	TotalTodos int64         `json:"totalTodos"`
}

func Index(c *gin.Context) {
	var todos []models.Todo
	var limit int = 10 // number of records per page
	var page int = 1   // default to the first page

	// parse the "page" query parameter
	if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		page = p
	}
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		limit = l
	}

	// calculate the offset based on the page number and limit
	offset := (page - 1) * limit

	// get user ID from token
	userID, err := getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// query the database for todos with the specified limit, offset, and user ID
	reqerror := initializers.DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&todos).Error

	// check for errors and set isError field if necessary
	if reqerror != nil {
		c.JSON(http.StatusInternalServerError, TodoIndexResponse{IsError: true})
		return
	}

	// count the total number of todos in the database
	var total int64
	initializers.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Count(&total)

	// calculate the total number of pages based on the total number of todos and the limit
	var totalPages int = int(math.Ceil(float64(total) / float64(limit)))

	// define the response struct and populate it with data
	response := TodoIndexResponse{
		IsError:    false,
		Todos:      todos,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalTodos: total,
	}

	// return the response as JSON
	c.JSON(http.StatusOK, response)
}

func Create(c *gin.Context) {
	var todo models.Todo

	// Parse user ID from JWT token
	userID, err := getUserIdFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"isError" : true,
			"message": "Unauthorized",
		})
		return
	}

	// Bind JSON request body to todo struct
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"isError" : true,
			"message": err.Error(),
		})
		return
	}

	// Set UserID field to the user ID parsed from the JWT token
	todo.UserID = uint(userID)

	// Create new todo in the database
	initializers.DB.Create(&todo)

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"isError" : false,
		"todo": todo,
	})
}

func getUserIdFromToken(c *gin.Context) (int64, error) {
	// Get token from Authorization header
	authHeader, err := c.Cookie("Authorization")
	if err != nil {
		return 0, errors.New("Authorization header missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return 0, fmt.Errorf("Error parsing JWT token: %v", err)
	}

	// Verify token is valid
	if !token.Valid {
		return 0, errors.New("Invalid JWT token")
	}

	// Get user ID from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Error getting user ID from JWT token")
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("Error getting user ID from JWT token")
	}

	userID := int64(userIDFloat)

	return userID, nil
}

func Update(c *gin.Context) {
	var todo models.Todo
	id := c.Param("id")

	if err := initializers.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"isError" : true,
				"message": "Todo not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"isError" : true,
			"message": "Failed to find Todo",
		})
		return
	}
	userID, err := getUserIdFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"isError" : true,
			"message": "Unauthorized",
		})
		return
	}
	// check if the user ID of the todo matches the user ID from the token
	if todo.UserID != uint(userID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"isError" : true,
			"message": "Unauthorized",
		})
		return
	}

	if err := initializers.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"isError" : true,
			"message": "Todo not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"isError" : true,
			"message": err.Error(),
		})
		return
	}

	initializers.DB.Save(&todo)
	c.JSON(http.StatusOK, gin.H{
		"isError" : false,
		"todos": todo,
	})
}

func Destroy(c *gin.Context) {
	var todo models.Todo
	id := c.Param("id")

	// find the todo by ID
	if err := initializers.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"isError" : true,
				"message": "Todo not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"isError" : true,
			"message": "Failed to find Todo",
		})
		return
	}
	userID, err := getUserIdFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"isError" : true,
			"message": "Unauthorized",
		})
		return
	}
	// check if the user ID of the todo matches the user ID from the token
	if todo.UserID != uint(userID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"isError" : true,
			"message": "Unauthorized",
		})
		return
	}

	// delete the todo
	if err := initializers.DB.Delete(&todo).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"isError" : true,
			"message": "Failed to delete Todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isError" : false,
		"message": "Todo deleted",
	})
}
