package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/igromancer/deck-assignment/internal/config"
	"github.com/igromancer/deck-assignment/internal/data"
)

const ctxAPIKeyID = "api_key_id"

func NewServer() (*Server, error) {
	cfg := config.GetConfig()
	jobRepo, err := data.NewJobrepository(*cfg)
	if err != nil {
		return nil, err
	}
	apiKeyRepo, err := data.NewApiKeyrepository(*cfg)
	if err != nil {
		return nil, err
	}
	s := &Server{
		router:     gin.Default(),
		apiKeyRepo: apiKeyRepo,
		jobRepo:    jobRepo,
		cfg:        *cfg,
	}

	return s, nil
}

type Server struct {
	router     *gin.Engine
	jobRepo    data.IJobRepository
	apiKeyRepo data.IApiKeyRepository
	cfg        config.Config
}

func (s *Server) Listen(addr ...string) error {
	s.router.POST("/apikey", s.createApiKey)

	authorized := s.router.Group("/")
	authorized.Use(AuthRequired(s.apiKeyRepo))
	authorized.POST("/jobs", s.createJob)
	authorized.GET("/jobs/:id", s.getJobStatus)
	authorized.GET("/jobs", s.listJobs)
	authorized.GET("/jobs/:id/result", s.getJobResult)
	return s.router.Run(addr...) // Default 0.0.0.0:8080
}

func (s *Server) createApiKey(c *gin.Context) {
	apiKey, err := s.apiKeyRepo.Create(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"apikey": apiKey})
}

func (s *Server) createJob(c *gin.Context) {
	body, err := ValidateCreateJobRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	apiKeyId, ok := apiKeyIDFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "api key was not set"})
		return
	}
	job := data.Job{
		Url:      body.Url,
		Status:   data.JobStatusPending,
		ApiKeyID: apiKeyId,
	}
	err = s.jobRepo.Create(c.Request.Context(), &job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO: Implement adding job to queue
	c.JSON(http.StatusAccepted, data.ToJobPublic(&job))
}

func (s *Server) listJobs(c *gin.Context) {
	apiKeyId, ok := apiKeyIDFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "api key was not set"})
		return
	}
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset provided"})
		return
	}
	limitStr := c.DefaultQuery("limit", "25")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit provided"})
		return
	}
	jobs, err := s.jobRepo.List(c.Request.Context(), apiKeyId, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	publicJobs := make([]data.JobPublic, len(jobs))
	for i := range jobs {
		publicJobs[i] = data.ToJobPublic(&jobs[i])
	}
	c.JSON(http.StatusOK, publicJobs)
}

func (s *Server) getJobStatus(c *gin.Context) {
	jobId := c.Param("id")
	uid, err := strconv.ParseUint(jobId, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}
	apiKeyId, ok := apiKeyIDFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "api key was not set"})
		return
	}
	job, err := s.jobRepo.Get(c.Request.Context(), uint(uid))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if apiKeyId != job.ApiKeyID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized to view this job status"})
		return
	}
	c.JSON(http.StatusOK, data.ToJobPublic(job))
}

func (s *Server) getJobResult(c *gin.Context) {
	jobId := c.Param("id")
	c.JSON(http.StatusNotImplemented, gin.H{
		"status": "Not implemented",
		"jobId":  jobId,
	})
}

func apiKeyIDFromContext(c *gin.Context) (uint, bool) {
	v, ok := c.Get(ctxAPIKeyID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
