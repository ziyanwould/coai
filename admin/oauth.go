package admin

import (
	"chat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthConfigForm struct {
	LinuxDo utils.LinuxDoOAuth `json:"linux_do"`
}

// GetOAuthConfigAPI 获取 OAuth 配置
func GetOAuthConfigAPI(c *gin.Context) {
	cfg := utils.GetOAuthConfig()
	// 不返回 client_secret，只返回是否已配置
	response := gin.H{
		"linux_do": gin.H{
			"enabled":      cfg.LinuxDo.Enabled,
			"client_id":    cfg.LinuxDo.ClientID,
			"redirect_url": cfg.LinuxDo.RedirectURL,
			"has_secret":   len(cfg.LinuxDo.ClientSecret) > 0,
		},
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   response,
	})
}

// UpdateOAuthConfigAPI 更新 OAuth 配置
func UpdateOAuthConfigAPI(c *gin.Context) {
	var form OAuthConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	// 如果没有提供新的 secret，保留旧的
	currentCfg := utils.GetOAuthConfig()
	if form.LinuxDo.ClientSecret == "" {
		form.LinuxDo.ClientSecret = currentCfg.LinuxDo.ClientSecret
	}

	err := utils.UpdateOAuthConfig(form.LinuxDo)

	c.JSON(http.StatusOK, gin.H{
		"status": err == nil,
		"error":  utils.GetError(err),
	})
}
