package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

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
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "Not implemented",
		"jobId":  jobId,
	})
}

func getJobResult(c *gin.Context) {
	jobId := c.Param("id")
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "Not implemented",
		"jobId":  jobId,
	})
}

func Listen(addr ...string) {
	router = gin.Default()

	router.POST("/jobs", createJob)
	router.GET("/jobs", listJobs)
	router.GET("/jobs/:id", getJobStatus)
	router.GET("/jobs/:id/result", getJobResult)
	router.Run(addr...) // Default 0.0.0.0:8080
}
