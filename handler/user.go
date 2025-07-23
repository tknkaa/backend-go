package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"myapp/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(c echo.Context) error {
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	newUser := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if result := h.DB.Create(&newUser); result.Error != nil {
		return c.String(http.StatusInternalServerError, "could not create user"+result.Error.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "signed up successfully",
	})
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignIn(c echo.Context) error {
	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var user model.User
	if result := h.DB.Where("userame = ?", req.Username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.String(http.StatusUnauthorized, "no user found")
		}
		return c.String(http.StatusInternalServerError, "database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.String(http.StatusUnauthorized, "invalid password")
	}

	sessionID := generateSessionID()
	newSession := model.Session{
		Token:     sessionID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserId:    user.ID,
	}

	if result := h.DB.Create(&newSession); result.Error != nil {
		return c.String(http.StatusInternalServerError, "could not create session")
	}

	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = sessionID
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "signed in successfully",
	})
}

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.String(http.StatusUnauthorized, "no session cookie")
		}
		sessionID := cookie.Value

		var session model.Session
		result := h.DB.Where("token = ?", sessionID).Preload("User").First(&session)

		if result.Error != nil {
			return c.String(http.StatusUnauthorized, "invalid session")
		}

		if session.ExpiresAt.Before(time.Now()) {
			h.DB.Delete(&session)
			return c.String(http.StatusUnauthorized, "session expired")
		}
		c.Set("user", session.User)

		return next(c)
	}
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
