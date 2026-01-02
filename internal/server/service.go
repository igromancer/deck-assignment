package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/igromancer/deck-assignment/internal/data"
)

var router *gin.Engine
var cfg *config.Config

func createApiKey(c *gin.Context) {
	repo, err := data.NewApiKeyrepository()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	apiKey, err := repo.Create(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"apikey": apiKey})
}

func createJob(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "Not implemented",
	})
}

func listJobs(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "Not implemented",
	})
}

func getJobStatus(c *gin.Context) {
	jobId := c.Param("id")
	repo, err := data.NewJobrepository()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uid, err := strconv.ParseUint(jobId, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}
	job, err := repo.Get(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func getJobResult(c *gin.Context) {
	jobId := c.Param("id")
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "Not implemented",
		"jobId":  jobId,
	})
}

func Listen(addr ...string) {
	cfg = config.GetConfig()
	router = gin.Default()

	router.POST("/apikey", createApiKey)
	router.POST("/jobs", createJob)
	router.GET("/jobs", listJobs)
	router.GET("/jobs/:id", getJobStatus)
	router.GET("/jobs/:id/result", getJobResult)
	router.Run(addr...) // Default 0.0.0.0:8080
}
