package financialhandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianwu291/go-learn/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service for testing
type mockFinancialService struct {
	mock.Mock
}

func (m *mockFinancialService) CalculateFinancial(req types.FinancialRawInfo, roundingDigits int) (types.FinancialResultInfo, error) {
	args := m.Called(req, roundingDigits)
	return args.Get(0).(types.FinancialResultInfo), args.Error(1)
}

func TestFinancialHandler_Calculate(t *testing.T) {
	// Setup Gin test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*mockFinancialService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Successful calculation",
			requestBody: types.FinancialRawInfo{
				Revenue:  1000,
				Expenses: 600,
				TaxRate:  0.2,
			},
			setupMock: func(m *mockFinancialService) {
				m.On("CalculateFinancial",
					types.FinancialRawInfo{Revenue: 1000, Expenses: 600, TaxRate: 0.2},
					2,
				).Return(types.FinancialResultInfo{Profit: 320.00, Ratio: 1.25}, nil)
			},
			expectedStatus: 200,
			expectedBody: &types.FinancialResultInfo{
				Profit: 320.00,
				Ratio:  1.25,
			},
		},
		{
			name: "Invalid request - negative revenue",
			requestBody: types.FinancialRawInfo{
				Revenue:  -1000,
				Expenses: 600,
				TaxRate:  0.2,
			},
			setupMock:      func(m *mockFinancialService) {}, // No mock needed for validation error
			expectedStatus: 400,
			expectedBody: &types.BadRequestResponse{
				Message: "Invalid request: Key: 'FinancialRawInfo.Revenue' Error:Field validation for 'Revenue' failed on the 'gte' tag",
			},
		},
		{
			name:           "Empty request body",
			requestBody:    nil,
			setupMock:      func(m *mockFinancialService) {}, // No mock needed for empty body
			expectedStatus: 400,
			expectedBody: &types.BadRequestResponse{
				Message: "Request body is empty",
			},
		},
		{
			name: "Invalid tax rate",
			requestBody: types.FinancialRawInfo{
				Revenue:  1000,
				Expenses: 600,
				TaxRate:  1.5, // Greater than 1
			},
			setupMock:      func(m *mockFinancialService) {}, // No mock needed for validation error
			expectedStatus: 400,
			expectedBody: &types.BadRequestResponse{
				Message: "Invalid request: Key: 'FinancialRawInfo.TaxRate' Error:Field validation for 'TaxRate' failed on the 'lte' tag",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mockFinancialService)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			handler := NewFinancialHandler(mockService)

			router := gin.New()
			router.POST("/calculate", handler.Calculate)

			// Create request
			var req *http.Request
			if tt.requestBody != nil {
				jsonBody, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBuffer(jsonBody))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/calculate", nil)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert response status
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			var response interface{}
			if tt.expectedStatus == 200 {
				response = &types.FinancialResultInfo{}
			} else {
				response = &types.BadRequestResponse{}
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verify that mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
