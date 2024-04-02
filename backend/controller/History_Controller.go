package controller

import (
	"audiohub/models"
	"audiohub/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveToHistory(c *gin.Context) {
	var history models.History
	if err := c.BindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	message, err := service.SaveToHistory(history)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}


func DisplayHistory(c *gin.Context) {
	var token models.Token
	if err := c.BindJSON(&token); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid  data"})
		return
	}
	result, message, err := service.DisplayHistory(token)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result})
}


func DeleteHistory(c *gin.Context) {
	var history models.DeteleHistory
	if err := c.BindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	message, err := service.DeteleHistory(history)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}
