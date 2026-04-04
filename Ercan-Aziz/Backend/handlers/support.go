package handlers

import (
	"context"
	"time"

	"notify-api/db"
	"notify-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostFeedback(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var feedback models.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(400, gin.H{"error": "Veri formatı hatalı! " + err.Error()})
		return
	}

	if feedback.Subject == "" || feedback.Message == "" {
		c.JSON(400, gin.H{"error": "Subject ve message alanları boş olamaz!"})
		return
	}
	if feedback.Type != "bug" && feedback.Type != "suggestion" && feedback.Type != "other" {
		c.JSON(400, gin.H{"error": "type alanı yalnızca 'bug', 'suggestion' veya 'other' olabilir!"})
		return
	}

	feedback.ID          = primitive.NewObjectID()
	feedback.UserID      = userID
	feedback.SubmittedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.FeedbackCollection.InsertOne(ctx, feedback)
	if err != nil {
		c.JSON(500, gin.H{"error": "Geri bildirim kaydedilemedi!"})
		return
	}

	c.JSON(201, gin.H{
		"message":      "Geri bildiriminiz alındı. En kısa sürede değerlendirilerek dönüş yapılacaktır.",
		"feedback_id":  feedback.ID,
		"submitted_at": feedback.SubmittedAt,
	})
}
