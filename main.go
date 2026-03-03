package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()

	r := gin.Default()

	// Enable CORS if needed (simple version)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve frontend
	r.StaticFile("/", "./static/index.html")
	r.Static("/static", "./static")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/devices", func(c *gin.Context) {
		devices, err := GetDevices()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, devices)
	})

	r.POST("/devices", func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
			MAC  string `json:"mac" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := AddDevice(input.Name, input.MAC)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "device added"})
	})

	r.DELETE("/devices/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		err = DeleteDevice(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "device deleted"})
	})

	r.POST("/wake/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		devices, err := GetDevices()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var targetMAC string
		for _, d := range devices {
			if d.ID == id {
				targetMAC = d.MAC
				break
			}
		}

		if targetMAC == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}

		err = WakeDevice(targetMAC)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "magic packet sent", "mac": targetMAC})
	})

	// Also support waking by MAC directly
	r.POST("/wake/mac/:mac", func(c *gin.Context) {
		mac := c.Param("mac")
		err := WakeDevice(mac)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "magic packet sent", "mac": mac})
	})

	log.Println("Server starting on :8090...")
	r.Run(":8090")
}
