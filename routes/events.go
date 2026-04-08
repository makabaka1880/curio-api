// Created by Sean L. on Apr. 8.
// Last Updated by Sean L. on Apr. 8.
//
// curio-api
// routes/events.go
//
// Makabaka1880, 2026. All rights reserved.

package routes

import (
	"curio-api/middleware"
	"curio-api/persistence"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleEventsRoutes(r *gin.Engine) *gin.RouterGroup {
	events := r.Group("/events")
	events.POST("/create", middleware.AuthMiddleware, HandleEventCreation)
	events.DELETE("/retract/:id", middleware.AuthMiddleware, HandleEventDeletion)
	events.GET("/list", HandleEventListing)
	return events
}

func HandleEventCreation(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		StartDate   string `json:"startDate"`
		EventDate   string `json:"eventDate"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON provided"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid start date format"})
		return
	}

	eventTime, err := time.Parse(time.RFC3339, input.EventDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid event date format"})
		return
	}

	eventObject := persistence.Event{
		Name:             input.Name,
		Description:      input.Description,
		SubmissionStarts: startTime,
		EventDate:        eventTime,
	}

	err = persistence.DB.Create(&eventObject).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create event", "spec": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Event created successfully",
		"data":    eventObject,
	})
}

func HandleEventDeletion(c *gin.Context) {
	id := c.Param("id")

	var entry persistence.Event
	err := persistence.DB.First(&entry, "id=?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Event not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error", "spec": err.Error()})
			return
		}
	}
	persistence.DB.Delete(&persistence.Event{}, "id=?", id)

	c.JSON(200, gin.H{
		"message": "Event deletion successful",
		"data": gin.H{
			"name": entry.Name,
		},
	})
}

func HandleEventListing(c *gin.Context) {
	filterStarts := c.Query("starts_after")
	filterEnds := c.Query("ends_before")
	pageSize := c.DefaultQuery("page_size", "10")
	pageNumber := c.DefaultQuery("page_number", "1")

	query := persistence.DB.Model(&persistence.Event{})

	if filterStarts != "" {
		if t, err := time.Parse(time.RFC3339, filterStarts); err == nil {
			query = query.Where("submission_start > ?", t)
		}
	}

	if filterEnds != "" {
		if t, err := time.Parse(time.RFC3339, filterEnds); err == nil {
			query = query.Where("submission_end < ?", t)
		}
	}

	ps, _ := strconv.Atoi(pageSize)
	pn, _ := strconv.Atoi(pageNumber)
	if ps <= 0 {
		ps = 10
	}
	if pn <= 0 {
		pn = 1
	}

	offset := (pn - 1) * ps

	var events []persistence.Event
	result := query.Offset(offset).Limit(ps).Find(&events)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error", "details": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Events retrieved successfully",
		"data":    events,
		"count":   len(events),
	})
}
