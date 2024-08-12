package user

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

var CLIENT_URL = os.Getenv("CLIENT_URL")

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

	c.JSON(http.StatusCreated, userRes)
}

func (h *Handler) Login(c *gin.Context) {
	refreshToken, err := c.Request.Cookie("refreshToken")

	/* Handle if user already logged in */
	if err == nil && refreshToken.Value != "" {
		isExist := h.Service.isSessionExist(c, RefreshToken(refreshToken.Value))
		if isExist {
			c.JSON(http.StatusConflict, gin.H{"message": "User already logged in"})
			return
		}

	}

	var reqUser LoginUserReq
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	userInfo := UserInfo{
		Phone:     reqUser.Phone,
		Password:  reqUser.Password,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	u, err := h.Service.Login(c.Request.Context(), &userInfo)
	if err != nil {

		if err.Error() == "invalid phone number" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshTokenMaxAge := int((30 * 24 * time.Hour).Seconds()) // 30 days

	c.SetCookie("refreshToken", u.RefreshToken, refreshTokenMaxAge, "/", CLIENT_URL, false, true)

	c.JSON(http.StatusOK, UserResponseWithAccess{
		BaseUserResponse: BaseUserResponse{
			ID:       int(u.ID),
			Username: u.Username,
			Phone:    u.Phone,
		},
		AccessToken: u.AccessToken,
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

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	us, err := h.Service.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, us)
}

func (h *Handler) GetUsersByIDs(c *gin.Context) {
	var usersIDs []int

	if err := c.ShouldBindJSON(&usersIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	us, err := h.Service.GetUsersByIDs(c, usersIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, us)
}

func (h *Handler) GetUsersByGroupID(c *gin.Context) {
	groupID := c.Query("group_id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Missing required query parameter group_id"})
		return
	}

	id, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "group_id is not a number"})
		return
	}

	users, err := h.Service.GetUsersByGroupID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Request.Cookie("refreshToken")
	if err != nil || refreshToken.Value == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	userWithTokens, err := h.Service.ValidateSession(c, RefreshToken(refreshToken.Value))
	if err != nil {
		if err.Error() == "Unauthorized" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshTokenMaxAge := int((30 * 24 * time.Hour).Seconds()) // 30 days

	c.SetCookie("refreshToken", string(userWithTokens.RefreshToken), refreshTokenMaxAge, "/", CLIENT_URL, false, true)

	c.JSON(http.StatusOK, UserResponseWithAccess{
		BaseUserResponse: BaseUserResponse{
			ID:       userWithTokens.ID,
			Username: userWithTokens.Username,
			Phone:    userWithTokens.Phone,
		},
		AccessToken: userWithTokens.AccessToken,
	})
}
