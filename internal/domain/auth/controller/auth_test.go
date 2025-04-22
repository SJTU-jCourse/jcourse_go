package controller

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jcourse_go/internal/domain/auth/model"
)

type mockVerificationCodeService struct {
	mock.Mock
}

func (m *mockVerificationCodeService) GenerateCode() string {
	return "123456"
}

func (m *mockVerificationCodeService) VerifyCode(ctx context.Context, email, code string) error {
	args := m.Called(ctx, email, code)
	return args.Error(0)
}

func (m *mockVerificationCodeService) SendCode(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (*model.UserDomain, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*model.UserDomain), args.Error(1)
}

func (m *mockAuthService) Register(ctx context.Context, email, password string) (*model.UserDomain, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*model.UserDomain), args.Error(1)
}

func (m *mockAuthService) ResetPassword(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

type mockEmailValidator struct {
	mock.Mock
}

func (m *mockEmailValidator) Validate(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func setupRouter(store sessions.Store, authController *AuthController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(sessions.Sessions("session", store))
	authGroup := r.Group("/")
	RegisterAuthController(authGroup, authController)
	gob.Register(&model.UserDomain{})
	return r
}

func TestAuthController_Login(t *testing.T) {
	mockAuthService := new(mockAuthService)
	mockVerificationService := new(mockVerificationCodeService)
	mockValidator := new(mockEmailValidator)

	controller := &AuthController{
		verificationCodeService: mockVerificationService,
		authService:             mockAuthService,
		emailValidator:          mockValidator,
	}

	store := sessions.NewCookieStore([]byte("secret"))
	router := setupRouter(store, controller)

	t.Run("success", func(t *testing.T) {
		user := &model.UserDomain{ID: 1, Email: "test@example.com"}
		mockAuthService.On("Login", mock.Anything, "test@example.com", "password").Return(user, nil)

		body := `{"email": "test@example.com", "password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("wrong params", func(t *testing.T) {
		body := `invalid_json`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("login error", func(t *testing.T) {
		mockAuthService.On("Login", mock.Anything, "user@example.com", "wrongpassword").Return(&model.UserDomain{}, errors.New("login error"))

		body := `{"email": "user@example.com", "password": "wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockAuthService.AssertExpectations(t)
	})
}

func TestAuthController_Logout(t *testing.T) {
	controller := &AuthController{}
	store := sessions.NewCookieStore([]byte("secret"))
	router := setupRouter(store, controller)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestAuthController_Register(t *testing.T) {
	mockAuthService := new(mockAuthService)
	mockVerificationService := new(mockVerificationCodeService)
	mockValidator := new(mockEmailValidator)

	controller := &AuthController{
		verificationCodeService: mockVerificationService,
		authService:             mockAuthService,
		emailValidator:          mockValidator,
	}

	store := sessions.NewCookieStore([]byte("secret"))
	router := setupRouter(store, controller)

	t.Run("success", func(t *testing.T) {
		mockValidator.On("Validate", "test@example.com").Return(true)
		mockVerificationService.On("VerifyCode", mock.Anything, "test@example.com", "1234").Return(nil)
		user := &model.UserDomain{ID: 1, Email: "test@example.com"}
		mockAuthService.On("Register", mock.Anything, "test@example.com", "password").Return(user, nil)

		body := `{"email": "test@example.com", "password": "password", "code": "1234"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		mockValidator.On("Validate", "invalid_email").Return(false)

		body := `{"email": "invalid_email", "password": "password", "code": "1234"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestAuthController_SendVerificationCode(t *testing.T) {
	mockVerificationService := new(mockVerificationCodeService)
	mockValidator := new(mockEmailValidator)

	controller := &AuthController{
		verificationCodeService: mockVerificationService,
		emailValidator:          mockValidator,
	}

	store := sessions.NewCookieStore([]byte("secret"))
	router := setupRouter(store, controller)

	t.Run("success", func(t *testing.T) {
		mockValidator.On("Validate", "test@example.com").Return(true)
		mockVerificationService.On("SendCode", mock.Anything, "test@example.com").Return(nil)

		body := `{"email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/send-verification-code", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockVerificationService.AssertExpectations(t)
	})
}

func TestAuthController_ResetPassword(t *testing.T) {
	mockAuthService := new(mockAuthService)
	mockValidator := new(mockEmailValidator)

	controller := &AuthController{
		authService:    mockAuthService,
		emailValidator: mockValidator,
	}

	store := sessions.NewCookieStore([]byte("secret"))
	router := setupRouter(store, controller)

	t.Run("success", func(t *testing.T) {
		mockValidator.On("Validate", "test@example.com").Return(true)
		mockAuthService.On("ResetPassword", mock.Anything, "test@example.com", "newpassword").Return(nil)

		body := `{"email": "test@example.com", "password": "newpassword", "code": "123456"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/reset-password", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockAuthService.AssertExpectations(t)
	})
}
