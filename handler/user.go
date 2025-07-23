package handler

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"myapp/model"
	"net/http"
	"time"
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

	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	sessionID := base64.URLEncoding.EncodeToString(randomBytes)

	newSession := model.Session{
		Token:     sessionID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserId:    newUser.ID,
	}

	if result := h.DB.Create(&newSession); result.Error != nil {
		return c.String(http.StatusInternalServerError, "could not create session"+result.Error.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = sessionID
	c.SetCookie(cookie)
	return c.JSON(http.StatusCreated, echo.Map{
		"id":        newUser.ID,
		"username":  newUser.Username,
		"createdAt": newUser.CreatedAt,
	})
}
