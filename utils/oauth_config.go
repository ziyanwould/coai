/*
 * @Author: Liu Jiarong
 * @Date: 2026-01-07
 * @Description: OAuth 独立配置管理，避免被 viper.WriteConfig() 覆盖
 */
package utils

import (
	"chat/globals"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// LinuxDoOAuth Linux.do OAuth 配置
type LinuxDoOAuth struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	ClientID     string `yaml:"client_id" json:"client_id"`
	ClientSecret string `yaml:"client_secret" json:"client_secret"`
	RedirectURL  string `yaml:"redirect_url" json:"redirect_url"`
}

// OAuthConfig OAuth 配置结构
type OAuthConfig struct {
	LinuxDo LinuxDoOAuth `yaml:"linux_do" json:"linux_do"`
}

var (
	oauthConfig     *OAuthConfig
	oauthConfigOnce sync.Once
	oauthConfigMu   sync.RWMutex
)

const oauthConfigFile = "config/oauth.yaml"

// LoadOAuthConfig 加载 OAuth 配置
func LoadOAuthConfig() *OAuthConfig {
	oauthConfigOnce.Do(func() {
		oauthConfig = &OAuthConfig{
			LinuxDo: LinuxDoOAuth{
				Enabled:      false,
				ClientID:     "",
				ClientSecret: "",
				RedirectURL:  "",
			},
		}
		reloadOAuthConfig()
	})
	return oauthConfig
}

// reloadOAuthConfig 重新加载配置文件
func reloadOAuthConfig() {
	oauthConfigMu.Lock()
	defer oauthConfigMu.Unlock()

	data, err := os.ReadFile(oauthConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			globals.Info("[oauth_config] config file not found, using defaults")
			// 创建默认配置文件
			SaveOAuthConfigInternal(oauthConfig)
		} else {
			globals.Warn(fmt.Sprintf("[oauth_config] failed to read config: %s", err.Error()))
		}
		return
	}

	var cfg OAuthConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		globals.Warn(fmt.Sprintf("[oauth_config] failed to parse config: %s", err.Error()))
		return
	}

	oauthConfig.LinuxDo = cfg.LinuxDo

	globals.Info(fmt.Sprintf("[oauth_config] loaded linux_do oauth config, enabled=%v",
		oauthConfig.LinuxDo.Enabled))
}

// GetOAuthConfig 获取 OAuth 配置（只读）
func GetOAuthConfig() OAuthConfig {
	oauthConfigMu.RLock()
	defer oauthConfigMu.RUnlock()

	if oauthConfig == nil {
		return OAuthConfig{}
	}

	// 返回副本
	return OAuthConfig{
		LinuxDo: oauthConfig.LinuxDo,
	}
}

// GetLinuxDoOAuth 获取 Linux.do OAuth 配置
func GetLinuxDoOAuth() LinuxDoOAuth {
	cfg := GetOAuthConfig()
	return cfg.LinuxDo
}

// IsLinuxDoOAuthEnabled 检查 Linux.do OAuth 是否启用
func IsLinuxDoOAuthEnabled() bool {
	cfg := GetLinuxDoOAuth()
	return cfg.Enabled && len(cfg.ClientID) > 0 && len(cfg.ClientSecret) > 0
}

// UpdateOAuthConfig 更新 OAuth 配置
func UpdateOAuthConfig(linuxDo LinuxDoOAuth) error {
	oauthConfigMu.Lock()
	defer oauthConfigMu.Unlock()

	if oauthConfig == nil {
		oauthConfig = &OAuthConfig{}
	}

	oauthConfig.LinuxDo = linuxDo

	return SaveOAuthConfigInternal(oauthConfig)
}

// SaveOAuthConfigInternal 内部保存函数
func SaveOAuthConfigInternal(cfg *OAuthConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal oauth config: %w", err)
	}

	// 添加注释头
	header := "# OAuth 配置\n# 此文件独立于主配置文件，不受 viper.WriteConfig() 影响\n\n"
	content := header + string(data)

	if err := os.WriteFile(oauthConfigFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write oauth config: %w", err)
	}

	globals.Info("[oauth_config] config saved successfully")
	return nil
}
