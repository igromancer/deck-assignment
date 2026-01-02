package server

import "github.com/gin-gonic/gin"

type CreateJobRequest struct {
	Url string `json:"url" binding:"required"`
}

func ValidateCreateJobRequest(c *gin.Context) (*CreateJobRequest, error) {
	var json CreateJobRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		return nil, err
	}
	return &json, nil
}
