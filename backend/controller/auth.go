package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/songquanpeng/go-api-starter/common"
	"github.com/songquanpeng/go-api-starter/common/client"
	"github.com/songquanpeng/go-api-starter/model"
)

var oauthClient *client.OAuthClient

func InitOAuth(cfg *client.OAuthConfig) {
	oauthClient = client.NewOAuthClient(cfg)
}

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (h *AuthController) GetAuthURL(c *gin.Context) {
	authURL, verifier, err := oauthClient.BuildAuthURL(oauthClient.GetRedirectURI())
	if err != nil {
		common.ServerError(c, "Failed to build auth URL")
		return
	}

	state := uuid.New().String()
	if err := saveCodeExchange(state, verifier, oauthClient.GetRedirectURI()); err != nil {
		common.ServerError(c, "Failed to save code exchange")
		return
	}

	common.Success(c, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *AuthController) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		common.ParamError(c, "Missing code or state parameter")
		return
	}

	codeExchange, err := getCodeExchange(state)
	if err != nil {
		common.ParamError(c, "Invalid state parameter")
		return
	}

	if codeExchange.Used {
		common.ParamError(c, "Code already used")
		return
	}

	tokenResp, err := oauthClient.ExchangeCodeForToken(code, codeExchange.CodeVerifier, codeExchange.RedirectURI)
	if err != nil {
		common.ServerError(c, "Failed to exchange code for token")
		return
	}

	markCodeExchangeUsed(state)

	userInfo, err := oauthClient.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		common.ServerError(c, "Failed to get user info")
		return
	}

	userModel := &model.User{
		UserID:   userInfo.UserID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		Nickname: userInfo.Nickname,
	}

	user, err := createOrUpdateUser(userModel)
	if err != nil {
		common.ServerError(c, "Failed to create or update user")
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

	if err := createSession(session); err != nil {
		common.ServerError(c, "Failed to create session")
		return
	}

	common.Success(c, gin.H{
		"session_id": session.SessionID,
		"user":       user,
		"expires_in": tokenResp.ExpiresIn,
	})
}

func (h *AuthController) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		common.Unauthorized(c, "User not authenticated")
		return
	}
	common.Success(c, user)
}

func (h *AuthController) RefreshToken(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		common.Unauthorized(c, "Session not found")
		return
	}

	sessionModel := session.(*model.Session)

	tokenResp, err := oauthClient.RefreshToken(sessionModel.RefreshToken)
	if err != nil {
		common.ServerError(c, "Failed to refresh token")
		return
	}

	sessionModel.AccessToken = tokenResp.AccessToken
	sessionModel.RefreshToken = tokenResp.RefreshToken

	if err := updateSession(sessionModel); err != nil {
		common.ServerError(c, "Failed to update session")
		return
	}

	common.Success(c, gin.H{
		"access_token":  tokenResp.AccessToken,
		"refresh_token": tokenResp.RefreshToken,
		"expires_in":    tokenResp.ExpiresIn,
	})
}

func (h *AuthController) Logout(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		common.SuccessWithMessage(c, "Logged out successfully")
		return
	}

	sessionModel := session.(*model.Session)

	if err := deleteSession(sessionModel.SessionID); err != nil {
		common.ServerError(c, "Failed to delete session")
		return
	}

	logoutURL := oauthClient.BuildLogoutURL("")
	common.Success(c, gin.H{
		"logout_url": logoutURL,
	})
}

func (h *AuthController) GetAPIKey(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		common.Unauthorized(c, "User not authenticated")
		return
	}

	userModel := user.(*model.User)

	if userModel.APIKey != "" {
		common.Success(c, gin.H{
			"api_key":        userModel.APIKey,
			"api_key_secret": userModel.APIKeySecret,
		})
		return
	}

	session, _ := c.Get("session")
	sessionModel := session.(*model.Session)

	apiKeyInfo, err := oauthClient.GetAPIKey(sessionModel.AccessToken)
	if err != nil {
		common.ServerError(c, "Failed to get API key")
		return
	}

	userModel.APIKey = apiKeyInfo.APIKey
	userModel.APIKeySecret = apiKeyInfo.APIKeySecret

	if err := updateUser(userModel); err != nil {
		common.ServerError(c, "Failed to update user API key")
		return
	}

	common.Success(c, gin.H{
		"api_key":        apiKeyInfo.APIKey,
		"api_key_secret": apiKeyInfo.APIKeySecret,
		"expire_time":    apiKeyInfo.ExpireTime,
	})
}

func saveCodeExchange(state, verifier, redirectURI string) error {
	exchange := &model.CodeExchange{
		State:        state,
		CodeVerifier: verifier,
		RedirectURI:  redirectURI,
		Used:         false,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	return common.DB.Create(exchange).Error
}

func getCodeExchange(state string) (*model.CodeExchange, error) {
	var exchange model.CodeExchange
	err := common.DB.Where("state = ? AND used = ? AND expires_at > ?", state, false, time.Now()).First(&exchange).Error
	return &exchange, err
}

func markCodeExchangeUsed(state string) {
	common.DB.Model(&model.CodeExchange{}).Where("state = ?", state).Update("used", true)
}

func createOrUpdateUser(userInfo *model.User) (*model.User, error) {
	var existingUser model.User
	err := common.DB.Where("user_id = ?", userInfo.UserID).First(&existingUser).Error
	if err != nil {
		if err := common.DB.Create(userInfo).Error; err != nil {
			return nil, err
		}
		return userInfo, nil
	}

	existingUser.Username = userInfo.Username
	existingUser.Email = userInfo.Email
	existingUser.Avatar = userInfo.Avatar
	existingUser.Nickname = userInfo.Nickname
	existingUser.LastLogin = time.Now()

	if err := common.DB.Save(&existingUser).Error; err != nil {
		return nil, err
	}

	return &existingUser, nil
}

func createSession(session *model.Session) error {
	return common.DB.Create(session).Error
}

func updateSession(session *model.Session) error {
	return common.DB.Save(session).Error
}

func deleteSession(sessionID string) error {
	return common.DB.Where("session_id = ?", sessionID).Delete(&model.Session{}).Error
}

func updateUser(user *model.User) error {
	return common.DB.Save(user).Error
}
