package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	segmentation_service "segmentationService"
	"segmentationService/src/service"
	mock_service "segmentationService/src/service/mocks"
	"testing"
)

func TestHandler_createUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user segmentation_service.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           segmentation_service.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: ` {"username":"Joe"}`,
			inputUser: segmentation_service.User{
				Username: "Joe",
			},
			mockBehavior: func(s *mock_service.MockUser, user segmentation_service.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "Service Failure",
			inputBody: ` {"username":"Joe"}`,
			inputUser: segmentation_service.User{
				Username: "Joe",
			},
			mockBehavior: func(s *mock_service.MockUser, user segmentation_service.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:                "Empty username field",
			inputBody:           ` {"username":}`,
			mockBehavior:        func(s *mock_service.MockUser, user segmentation_service.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid username field"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputUser)

			services := &service.Service{User: user}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/create", handler.createUser)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(recorder, req)
			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestHandler_deleteUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user segmentation_service.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           segmentation_service.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: ` {"username":"Joe"}`,
			inputUser: segmentation_service.User{
				Username: "Joe",
			},
			mockBehavior: func(s *mock_service.MockUser, user segmentation_service.User) {
				s.EXPECT().DeleteUser(user).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"user was successful deleted"}`,
		},
		{
			name:      "Service Failure",
			inputBody: ` {"username":"Joe"}`,
			inputUser: segmentation_service.User{
				Username: "Joe",
			},
			mockBehavior: func(s *mock_service.MockUser, user segmentation_service.User) {
				s.EXPECT().DeleteUser(user).Return(errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:                "Empty username field",
			inputBody:           ` {"username":}`,
			mockBehavior:        func(s *mock_service.MockUser, user segmentation_service.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid username field"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.inputUser)

			services := &service.Service{User: user}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/delete", handler.deleteUser)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(recorder, req)
			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}
