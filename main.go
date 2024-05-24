package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	filePath = "notification.json"
)

type Notification struct {
	ID                 uint    `json:"id"`
	CurrentPrice       float64 `json:"current_price"`
	DailyChangePercent float64 `json:"daily_change_percent"`
	TradingVolume      float64 `json:"trading_volume"`
	Status             string  `json:"status"`
}

func main() {
	r := gin.Default()

	r.POST("/notifications", createNotification)
	r.GET("/notifications", listNotification)
	r.PUT("/notifications/:id", updateNotification)
	r.DELETE("/notifications/:id", deleteNotification)
	r.POST("/notifications/send/:id", sendNotification)

	r.Run(":8080")
}

func createNotification(c *gin.Context) {
	var notification Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notifications, err := loadNotifications()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification.ID = uint(len(notifications) + 1)
	notification.Status = "Pending"
	notifications = append(notifications, notification)
	if err := saveNotifications(notifications); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, notification)

}

func listNotification(c *gin.Context) {
	notifications, err := loadNotifications()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)

}

func deleteNotification(c *gin.Context) {
	id := c.Param("id")
	notifications, err := loadNotifications()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range notifications {
		if id == fmt.Sprintf("%d", notifications[i].ID) {
			notifications = append(notifications[:i], notifications[i+1:]...)
			if err := saveNotifications(notifications); err != nil {
				c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
}

func updateNotification(c *gin.Context) {
	id := c.Param("id")
	notifications, err := loadNotifications()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newNotification Notification
	if err := c.ShouldBindJSON(&newNotification); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range notifications {
		if id == fmt.Sprintf("%d", notifications[i].ID) {
			notifications[i] = newNotification
			if err := saveNotifications(notifications); err != nil {
				c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.Status(http.StatusOK)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
}

func sendNotification(c *gin.Context) {
	id := c.Param("id")
	notifications, err := loadNotifications()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range notifications {
		if id == fmt.Sprintf("%d", notifications[i].ID) {
			notifications[i].Status = "Sent"
			if err := saveNotifications(notifications); err != nil {
				c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.Status(http.StatusOK)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
}

func loadNotifications() ([]Notification, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Notification{}, nil
		}
		return nil, err
	}
	var notifications []Notification
	if err := json.Unmarshal(file, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil

}
func saveNotifications(notifications []Notification) error {
	file, err := json.MarshalIndent(notifications, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, file, 0644)
}
