package config

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"math"
	"net/http/httptest"
	"testing"
)

func TestGlobalResponseError(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	newErr := errors.New("test error")
	_ = GlobalResponse(nil, newErr, c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["data"])
	assert.Equal(t, newErr.Error(), response["error"])
	assert.Equal(t, float64(-1), response["status"])
}

func TestGlobalResponseSuccess(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	mockData := "success"
	_ = GlobalResponse(mockData, nil, c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["error"])
	assert.Equal(t, mockData, response["data"])
	assert.Equal(t, float64(1), response["status"])
}

func TestGlobalResponseSuccessFloatTrunk(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	mockData := 1.123456789101112131415
	trunkMockData := math.Floor(mockData*1e8) / 1e8
	_ = GlobalResponse(mockData, nil, c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["error"])
	assert.Equal(t, trunkMockData, response["data"])
	assert.Equal(t, float64(1), response["status"])
}
