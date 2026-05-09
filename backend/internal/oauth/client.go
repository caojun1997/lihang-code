package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"pdf-parser/internal/config"
)

type OAuthClient struct {
	config     *OAuthConfig
	httpClient *http.Client
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	RedirectURI  string
	APIKeyURL    string
	LogoutURL    string
	WellKnownURL string
}

type AuthCodeRequest struct {
	ClientID        string
	ResponseType    string
	RedirectURI     string
	CodeChallenge   string
	CodeChallengeMethod string
	State           string
}

type TokenRequest struct {
	ClientID       string
	ClientSecret   string
	GrantType      string
	Code           string
	CodeVerifier   string
	RedirectURI    string
	RefreshToken  string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type UserInfo struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
}

type APIKeyInfo struct {
	APIKey     string `json:"api_key"`
	APIKeySecret string `json:"api_key_secret"`
	ExpireTime string `json:"expire_time"`
}

func NewOAuthClient(cfg *config.OAuthConfig) *OAuthClient {
	return &OAuthClient{
		config: &OAuthConfig{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			AuthURL:      cfg.AuthURL,
			TokenURL:     cfg.TokenURL,
			UserInfoURL:  cfg.UserInfoURL,
			RedirectURI:  cfg.RedirectURI,
			APIKeyURL:    cfg.APIKeyURL,
			LogoutURL:    cfg.LogoutURL,
			WellKnownURL: cfg.WellKnownURL,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func GenerateCodeVerifier() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	const length = 128

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range result {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func GenerateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	d := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(d)
}

func (c *OAuthClient) BuildAuthURL(redirectURI string) (string, string, error) {
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	challenge := GenerateCodeChallenge(verifier)

	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("code_challenge", challenge)
	params.Set("code_challenge_method", "S256")
	params.Set("state", fmt.Sprintf("%d", time.Now().UnixNano()))

	authURL := fmt.Sprintf("%s?%s", c.config.AuthURL, params.Encode())

	return authURL, verifier, nil
}

func (c *OAuthClient) ExchangeCodeForToken(code, codeVerifier, redirectURI string) (*TokenResponse, error) {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("client_secret", c.config.ClientSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("code_verifier", codeVerifier)
	params.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("GET", c.config.TokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = params.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}

func (c *OAuthClient) RefreshToken(refreshToken string) (*TokenResponse, error) {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("client_secret", c.config.ClientSecret)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("GET", c.config.TokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = params.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}

func (c *OAuthClient) GetUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", c.config.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user info failed with status: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &userInfo, nil
}

func (c *OAuthClient) GetAPIKey(accessToken string) (*APIKeyInfo, error) {
	req, err := http.NewRequest("GET", c.config.APIKeyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get API key failed with status: %d", resp.StatusCode)
	}

	var apiKeyInfo APIKeyInfo
	if err := json.NewDecoder(resp.Body).Decode(&apiKeyInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiKeyInfo, nil
}

func (c *OAuthClient) GetClientToken() (*TokenResponse, error) {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("client_secret", c.config.ClientSecret)
	params.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("GET", c.config.TokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = params.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}

func (c *OAuthClient) BuildLogoutURL(redirectURI string) string {
	if redirectURI == "" {
		return c.config.LogoutURL
	}
	return fmt.Sprintf("%s?redirect_uri=%s", c.config.LogoutURL, url.QueryEscape(redirectURI))
}

func (c *OAuthClient) GetRedirectURI() string {
	return c.config.RedirectURI
}

func (c *OAuthClient) GetConfig() *OAuthConfig {
	return c.config
}
