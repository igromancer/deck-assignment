package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateJobRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	out, err := ValidateCreateJobRequest(c)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "https://example.com", out.Url)
}

func TestValidateCreateJobRequest_MissingURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	out, err := ValidateCreateJobRequest(c)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestValidateCreateJobRequest_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	out, err := ValidateCreateJobRequest(c)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestValidateCreateJobRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{invalid-json}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	out, err := ValidateCreateJobRequest(c)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestValidateCreateJobRequest_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"url":123}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	out, err := ValidateCreateJobRequest(c)
	require.Error(t, err)
	assert.Nil(t, out)
}
