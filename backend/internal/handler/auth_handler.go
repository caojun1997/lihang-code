package handler

import (
	"net/http"
	"time"

	"pdf-parser/internal/model"
	"pdf-parser/internal/oauth"
	"pdf-parser/internal/repository"
	"pdf-parser/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	oauthClient *oauth.OAuthClient
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthHandler(oauthClient *oauth.OAuthClient) *AuthHandler {
	return &AuthHandler{
		oauthClient: oauthClient,
		userRepo:    repository.NewUserRepository(),
		sessionRepo: repository.NewSessionRepository(),
	}
}

func (h *AuthHandler) GetAuthURL(c *gin.Context) {
	authURL, verifier, err := h.oauthClient.BuildAuthURL(h.oauthClient.GetRedirectURI())
	if err != nil {
		response.ServerError(c, "Failed to build auth URL")
		return
	}

	state := uuid.New().String()
	if err := h.sessionRepo.SaveCodeExchange(state, verifier, h.oauthClient.GetRedirectURI()); err != nil {
		response.ServerError(c, "Failed to save code exchange")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"auth_url": authURL,
			"state":    state,
		},
	})
}

func (h *AuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		response.ParamError(c, "Missing code or state parameter")
		return
	}

	codeExchange, err := h.sessionRepo.GetCodeExchange(state)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid state parameter")
		return
	}

	if codeExchange.Used {
		response.Error(c, http.StatusBadRequest, 400, "Code already used")
		return
	}

	tokenResp, err := h.oauthClient.ExchangeCodeForToken(code, codeExchange.CodeVerifier, codeExchange.RedirectURI)
	if err != nil {
		response.ServerError(c, "Failed to exchange code for token")
		return
	}

	h.sessionRepo.MarkCodeExchangeUsed(state)

	userInfo, err := h.oauthClient.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		response.ServerError(c, "Failed to get user info")
		return
	}

	userModel := &model.User{
		UserID:   userInfo.UserID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		Nickname: userInfo.Nickname,
	}

	user, err := h.userRepo.CreateOrUpdateUser(userModel)
	if err != nil {
		response.ServerError(c, "Failed to create or update user")
		return
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	session := &model.Session{
		SessionID:    uuid.New().String(),
		UserID:       user.UserID,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	}

	if err := h.sessionRepo.CreateSession(session); err != nil {
		response.ServerError(c, "Failed to create session")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"session_id": session.SessionID,
			"user":       user,
			"expires_in": tokenResp.ExpiresIn,
		},
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "User not authenticated")
		return
	}

	response.Success(c, user)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "Session not found")
		return
	}

	sessionModel := session.(*model.Session)

	tokenResp, err := h.oauthClient.RefreshToken(sessionModel.RefreshToken)
	if err != nil {
		response.ServerError(c, "Failed to refresh token")
		return
	}

	sessionModel.AccessToken = tokenResp.AccessToken
	sessionModel.RefreshToken = tokenResp.RefreshToken

	if err := h.sessionRepo.UpdateSession(sessionModel); err != nil {
		response.ServerError(c, "Failed to update session")
		return
	}

	response.Success(c, gin.H{
		"access_token":  tokenResp.AccessToken,
		"refresh_token": tokenResp.RefreshToken,
		"expires_in":    tokenResp.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		response.SuccessWithMessage(c, "Logged out successfully")
		return
	}

	sessionModel := session.(*model.Session)

	if err := h.sessionRepo.DeleteSession(sessionModel.SessionID); err != nil {
		response.ServerError(c, "Failed to delete session")
		return
	}

	logoutURL := h.oauthClient.BuildLogoutURL("")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"logout_url": logoutURL,
		},
	})
}

func (h *AuthHandler) GetAPIKey(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "User not authenticated")
		return
	}

	userModel := user.(*model.User)

	if userModel.APIKey != "" {
		response.Success(c, gin.H{
			"api_key":      userModel.APIKey,
			"api_key_secret": userModel.APIKeySecret,
		})
		return
	}

	session, _ := c.Get("session")
	sessionModel := session.(*model.Session)

	apiKeyInfo, err := h.oauthClient.GetAPIKey(sessionModel.AccessToken)
	if err != nil {
		response.ServerError(c, "Failed to get API key")
		return
	}

	userModel.APIKey = apiKeyInfo.APIKey
	userModel.APIKeySecret = apiKeyInfo.APIKeySecret

	if err := h.userRepo.UpdateUser(userModel); err != nil {
		response.ServerError(c, "Failed to update user API key")
		return
	}

	response.Success(c, gin.H{
		"api_key":        apiKeyInfo.APIKey,
		"api_key_secret": apiKeyInfo.APIKeySecret,
		"expire_time":    apiKeyInfo.ExpireTime,
	})
}
