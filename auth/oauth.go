package auth

import (
	"chat/channel"
	"chat/globals"
	"chat/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/oauth2"
)

// LinuxDoUserInfo Linux.do 用户信息结构
type LinuxDoUserInfo struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// getLinuxDoOAuthConfig 获取 Linux.do OAuth 配置
func getLinuxDoOAuthConfig() *oauth2.Config {
	var clientID, clientSecret, redirectURL string
	if channel.SystemInstance != nil {
		state := channel.SystemInstance.OAuth.LinuxDo
		clientID = state.ClientID
		clientSecret = state.ClientSecret
		redirectURL = state.RedirectURL
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://connect.linux.do/oauth2/authorize",
			TokenURL: "https://connect.linux.do/oauth2/token",
		},
	}
}

// isLinuxDoOAuthEnabled 检查 Linux.do OAuth 是否启用
func isLinuxDoOAuthEnabled() bool {
	if channel.SystemInstance == nil {
		return false
	}

	state := channel.SystemInstance.OAuth.LinuxDo

	return state.Enabled && len(state.ClientID) > 0 && len(state.ClientSecret) > 0
}

// generateOAuthState 生成 OAuth state 参数并存储到 Redis
func generateOAuthState(c *gin.Context, cache *redis.Client) string {
	state := utils.GenerateCode(32)
	cache.Set(c, fmt.Sprintf("nio:oauth:state:%s", state), "1", 10*time.Minute)
	return state
}

// validateOAuthState 验证 OAuth state 参数
func validateOAuthState(c *gin.Context, cache *redis.Client, state string) bool {
	key := fmt.Sprintf("nio:oauth:state:%s", state)
	val, err := cache.Get(c, key).Result()
	if err != nil || val != "1" {
		return false
	}

	// 验证后立即删除,防止重放攻击
	cache.Del(c, key)
	return true
}

// getLinuxDoUserInfo 获取 Linux.do 用户信息
func getLinuxDoUserInfo(ctx context.Context, token *oauth2.Token) (*LinuxDoUserInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://connect.linux.do/api/user", nil)
	if err != nil {
		globals.Warn(fmt.Sprintf("[oauth] failed to create request: %s", err.Error()))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	globals.Info(fmt.Sprintf("[oauth] requesting user info from: https://connect.linux.do/api/user"))

	resp, err := client.Do(req)
	if err != nil {
		globals.Warn(fmt.Sprintf("[oauth] failed to send request: %s", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	globals.Info(fmt.Sprintf("[oauth] user info response status: %d, body: %s", resp.StatusCode, string(body)))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s (status: %d)", string(body), resp.StatusCode)
	}

	var userInfo LinuxDoUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		globals.Warn(fmt.Sprintf("[oauth] failed to decode user info: %s", err.Error()))
		return nil, err
	}

	globals.Info(fmt.Sprintf("[oauth] successfully got user info: id=%d, username=%s, email=%s", userInfo.ID, userInfo.Username, userInfo.Email))

	return &userInfo, nil
}

// generateUniqueUsername 生成唯一的用户名
func generateUniqueUsername(db *sql.DB, baseUsername string) string {
	username := fmt.Sprintf("linux_do_%s", baseUsername)

	// 如果用户名不存在,直接返回
	if !IsUserExist(db, username) {
		return username
	}

	// 用户名已存在,添加随机后缀
	for i := 0; i < 100; i++ {
		randomSuffix := utils.GenerateCode(4)
		newUsername := fmt.Sprintf("%s_%s", username, randomSuffix)
		if !IsUserExist(db, newUsername) {
			return newUsername
		}
	}

	// 极端情况:添加时间戳
	return fmt.Sprintf("%s_%d", username, time.Now().Unix())
}

// LinuxDoOAuthLogin 处理 Linux.do OAuth 登录
func LinuxDoOAuthLogin(c *gin.Context, userInfo *LinuxDoUserInfo) (string, error) {
	db := utils.GetDBFromContext(c)

	// 检查邮箱是否已存在
	if IsEmailExist(db, userInfo.Email) {
		// 邮箱已存在,直接登录
		user := GetUserByEmail(db, userInfo.Email)
		if user == nil {
			return "", errors.New("无法获取用户信息")
		}

		if user.IsBanned(db) {
			return "", errors.New("当前用户已被封禁")
		}

		// 更新 token 字段为 Linux.do ID
		linuxDoToken := fmt.Sprintf("linux_do:%d", userInfo.ID)
		globals.QueryDb(db, "UPDATE auth SET token = ? WHERE id = ?", linuxDoToken, user.ID)

		return user.GenerateTokenSafe(db)
	}

	// 新用户,自动注册
	if globals.CloseRegistration {
		return "", errors.New("本站暂未开放注册")
	}

	// 生成唯一用户名
	username := generateUniqueUsername(db, userInfo.Username)

	// 生成随机密码
	password := utils.GenerateChar(64)
	hash := utils.Sha2Encrypt(password)

	// 创建用户
	linuxDoToken := fmt.Sprintf("linux_do:%d", userInfo.ID)

	user := &User{
		Username: username,
		Password: hash,
		Email:    userInfo.Email,
		BindID:   getMaxBindId(db) + 1,
		Token:    linuxDoToken,
	}

	if _, err := globals.ExecDb(db, `
		INSERT INTO auth (username, password, email, bind_id, token)
		VALUES (?, ?, ?, ?, ?)
	`, user.Username, user.Password, user.Email, user.BindID, user.Token); err != nil {
		return "", fmt.Errorf("创建用户失败: %s", err.Error())
	}

	// 创建初始配额
	user.CreateInitialQuota(db)

	globals.Info(fmt.Sprintf("[oauth] new user registered via linux.do: %s (email: %s)", username, userInfo.Email))

	return user.GenerateToken()
}

// LinuxDoOAuthLoginAPI OAuth 授权跳转 API
func LinuxDoOAuthLoginAPI(c *gin.Context) {
	if !isLinuxDoOAuthEnabled() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Linux.do OAuth 未启用",
		})
		return
	}

	cache := utils.GetCacheFromContext(c)
	config := getLinuxDoOAuthConfig()

	// 生成 state 参数防止 CSRF 攻击
	state := generateOAuthState(c, cache)

	// 生成授权 URL
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// 检查是否需要返回 JSON (API 调用) 还是直接重定向 (浏览器调用)
	acceptHeader := c.GetHeader("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		// 前端 API 调用,返回 JSON
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"url":    authURL,
		})
	} else {
		// 直接重定向到 Linux.do 授权页面
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

// LinuxDoOAuthCallbackAPI OAuth 回调处理 API
func LinuxDoOAuthCallbackAPI(c *gin.Context) {
	if !isLinuxDoOAuthEnabled() {
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=oauth_disabled")
		return
	}

	cache := utils.GetCacheFromContext(c)

	// 获取参数
	code := c.Query("code")
	state := c.Query("state")
	errorMsg := c.Query("error")

	// 检查错误
	if len(errorMsg) > 0 {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/login?error=%s", errorMsg))
		return
	}

	// 验证 code
	if len(code) == 0 {
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=missing_code")
		return
	}

	// 验证 state
	if !validateOAuthState(c, cache, state) {
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=invalid_state")
		return
	}

	// 交换 token
	config := getLinuxDoOAuthConfig()
	ctx := context.Background()

	token, err := config.Exchange(ctx, code)
	if err != nil {
		globals.Warn(fmt.Sprintf("[oauth] failed to exchange token: %s", err.Error()))
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=exchange_failed")
		return
	}

	// 获取用户信息
	userInfo, err := getLinuxDoUserInfo(ctx, token)
	if err != nil {
		globals.Warn(fmt.Sprintf("[oauth] failed to get user info: %s", err.Error()))
		c.Redirect(http.StatusTemporaryRedirect, "/login?error=userinfo_failed")
		return
	}

	// 处理登录/注册
	jwtToken, err := LinuxDoOAuthLogin(c, userInfo)
	if err != nil {
		globals.Warn(fmt.Sprintf("[oauth] login failed: %s", err.Error()))
		errorEncoded := strings.ReplaceAll(err.Error(), " ", "_")
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/login?error=%s", errorEncoded))
		return
	}

	// 登录成功,重定向到前端并携带 token
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/oauth-success?token=%s", jwtToken))
}
