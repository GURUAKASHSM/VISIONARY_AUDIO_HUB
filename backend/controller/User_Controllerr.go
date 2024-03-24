package controller

import (
	"audiohub/constants"
	"audiohub/models"
	"audiohub/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Signup Function
func CreateProfile(c *gin.Context) {
	var profile models.Customer
	if err := c.BindJSON(&profile); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	result := service.CreateCustomer(profile)
	c.JSON(http.StatusOK, result)
}

// Customer Email Verification

func VerifyEmail(c *gin.Context) {
	var Data models.VerifyEmail
	if err := c.BindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	data, err := service.EmailVerification(Data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": data})
}

// Signin Function
func Login(c *gin.Context) {
	var request models.Login
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	token, no, err := service.Login(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if no == 1 {
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	} else if no == 0 {
		c.JSON(http.StatusOK, gin.H{"message": token})
		return
	}
}

// Validate Customer Token
func ValidateToken(c *gin.Context) {
	var userdata models.Token
	if err := c.BindJSON(&userdata); err != nil {

		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	result := service.Validatetoken(userdata.Token, constants.SecretKey)
	c.JSON(http.StatusOK, gin.H{"result": result})
}
func ForgetPassword(c *gin.Context) {
	var email models.ForgetPassword
	if err := c.BindJSON(&email); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid  data"})
		return
	}
	result, err := service.SendEmailForForgotPassword(email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": result})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result})
}
func PasswordChange(c *gin.Context) {
	var data models.PasswordChange
	if err := c.BindJSON(&data); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid  data"})
		return
	}
	result, err := service.ChangePassword(data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": result})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result})

}
func LisUserDetails(c *gin.Context) {
	var token models.Token
	if err := c.BindJSON(&token); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid  data"})
		return
	}
	result, message, _ := service.UserDetails(token)
	log.Print(message)
	c.JSON(http.StatusOK, gin.H{"result": result})

}
