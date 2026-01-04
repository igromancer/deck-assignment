package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/igromancer/deck-assignment/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Minimal mock implementing data.IApiKeyRepository
type MockApiKeyRepo struct{ mock.Mock }

func (m *MockApiKeyRepo) Get(ctx context.Context, publicId string) (*data.ApiKey, error) {
	args := m.Called(ctx, publicId)
	key, _ := args.Get(0).(*data.ApiKey)
	return key, args.Error(1)
}

func (m *MockApiKeyRepo) Create(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func setupAuthTestRouter(repo data.IApiKeyRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthRequired(repo))
	r.GET("/protected", func(c *gin.Context) {
		// ensure middleware sets api_key_id
		v, ok := c.Get("api_key_id")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "api_key_id missing"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"api_key_id": v})
	})
	return r
}

func TestAuthRequired_MissingHeader(t *testing.T) {
	repo := new(MockApiKeyRepo)
	r := setupAuthTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "x-api-key is required")
	repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

func TestAuthRequired_BadFormat(t *testing.T) {
	repo := new(MockApiKeyRepo)
	r := setupAuthTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("x-api-key", "no-dot-here")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "formatted incorrectly")
	repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

func TestAuthRequired_RepoError(t *testing.T) {
	repo := new(MockApiKeyRepo)
	r := setupAuthTestRouter(repo)

	// header has publicId.secret
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("x-api-key", "public.secret")
	w := httptest.NewRecorder()

	repo.On("Get", mock.Anything, "public").Return((*data.ApiKey)(nil), assert.AnError).Once()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), assert.AnError.Error())
	repo.AssertExpectations(t)
}

func TestAuthRequired_InvalidSecret(t *testing.T) {
	repo := new(MockApiKeyRepo)
	r := setupAuthTestRouter(repo)

	// Make a stored bcrypt hash for some other secret
	hashed, err := bcrypt.GenerateFromPassword([]byte("real-secret"), bcrypt.DefaultCost)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("x-api-key", "public.wrong-secret")
	w := httptest.NewRecorder()

	repo.On("Get", mock.Anything, "public").Return(&data.ApiKey{
		Model:        gormModelIDOnly(123),
		PublicId:     "public",
		HashedSecret: string(hashed),
	}, nil).Once()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid x-api-key")
	repo.AssertExpectations(t)
}

func TestAuthRequired_Success_SetsAPIKeyID(t *testing.T) {
	repo := new(MockApiKeyRepo)
	r := setupAuthTestRouter(repo)

	hashed, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("x-api-key", "public.secret")
	w := httptest.NewRecorder()

	repo.On("Get", mock.Anything, "public").Return(&data.ApiKey{
		Model:        gormModelIDOnly(777),
		PublicId:     "public",
		HashedSecret: string(hashed),
	}, nil).Once()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"api_key_id":777`)
	repo.AssertExpectations(t)
}

func gormModelIDOnly(id uint) gorm.Model {
	return gorm.Model{ID: id}
}
