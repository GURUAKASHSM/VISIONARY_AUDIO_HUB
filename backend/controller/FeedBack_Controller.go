package controller

import (
	"audiohub/models"
	"audiohub/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertFeedback(c *gin.Context) {
	var feedback models.Feedback

	if err := c.BindJSON(&feedback); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	result, err := service.InstertFeedback(feedback)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": result})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result})
}

// Delete FeedBack
func Deletefeedback(c *gin.Context) {
	var feedback models.FeedbackDB

	if err := c.BindJSON(&feedback); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	result := service.Deletefeedback(feedback)
	c.JSON(http.StatusOK, result)

}

// Get All FeedBacks
func GetFeedback(c *gin.Context) {
	data := service.GetFeedBacks()
	c.JSON(http.StatusOK, gin.H{"result": data})

}
