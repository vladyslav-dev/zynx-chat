package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var u CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRes, err := h.Service.Register(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userRes)
}

func (h *Handler) Login(c *gin.Context) {
	var reqUser LoginUserReq
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	userInfo := UserInfo{
		Email:     reqUser.Email,
		Password:  reqUser.Password,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	u, err := h.Service.Login(c.Request.Context(), &userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshTokenMaxAge := int((30 * 24 * time.Hour).Seconds()) // 30 days

	c.SetCookie("refreshToken", u.RefreshToken, refreshTokenMaxAge, "/", "localhost", false, true)

	c.JSON(http.StatusOK, UserRes{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Request.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.Logout(c, RefreshToken(refreshToken.Value))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successfully"})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	us, err := h.Service.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, us)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Request.Cookie("refreshToken")
	if err != nil || refreshToken.Value == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	tokens, err := h.Service.ValidateSession(c, RefreshToken(refreshToken.Value))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshTokenMaxAge := int((30 * 24 * time.Hour).Seconds()) // 30 days

	c.SetCookie("refreshToken", string(tokens.RefreshToken), refreshTokenMaxAge, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"accessToken": tokens.AccessToken})
}
