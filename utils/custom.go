/*
 * @Author: 刘家荣 14731753+liujiarong2@user.noreply.gitee.com
 * @Date: 2025-01-21 18:04:14
 * @LastEditors: Liu Jiarong
 * @LastEditTime: 2025-03-03 20:58:06
 * @FilePath: /coai/utils/custom.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koroFileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package utils

import (
	"chat/globals"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// Config 定义了配置文件的结构
type Config struct {
	CustomVisionModels []string `json:"custom_vision_models"`
}

var (
	customVisionConfig   Config
	lastConfigUpdateTime time.Time
	configMutex          sync.RWMutex
	remoteConfigURL      = "http://192.168.31.135/vision_models.json" // 远程配置文件的 URL
)

func init() {
	loadConfig() // 首次加载配置
	go func() {
		// 定时从本地文件重新加载配置 (每 5 分钟)
		for {
			time.Sleep(5 * time.Minute)
			loadLocalConfig() // 仅重新加载本地文件，并合并远程配置
		}
	}()
	go func() {
		// 定时从远程 URL 重新加载配置并合并 (每 1 小时)
		for {
			time.Sleep(1 * time.Hour)
			loadRemoteConfigAndMerge()
		}
	}()
}

// loadLocalConfig 仅加载本地配置文件，并与已有的远程配置合并
func loadLocalConfig() {
	localConfig, err := loadConfigFromFile("config/vision_models.json")
	if err != nil {
		globals.Warn(fmt.Sprintf("Error loading local config: %v", err))
		return
	}

	configMutex.Lock()
	defer configMutex.Unlock()
	remoteModels := customVisionConfig.CustomVisionModels                                                  // 保留当前配置中的远程模型
	customVisionConfig = *localConfig                                                                      // 更新为本地配置
	customVisionConfig.CustomVisionModels = append(customVisionConfig.CustomVisionModels, remoteModels...) // 重新合并远程模型
	customVisionConfig.CustomVisionModels = uniqueStringSlice(customVisionConfig.CustomVisionModels)       // 去重

	lastConfigUpdateTime = time.Now()
	globals.Info(fmt.Sprintf("Successfully reload local config, current config:%v, update time: %v", customVisionConfig, lastConfigUpdateTime))
}

// loadRemoteConfigAndMerge 加载远程配置文件并与本地配置合并
func loadRemoteConfigAndMerge() {
	remoteConfig, err := loadConfigFromURL(remoteConfigURL)
	if err != nil {
		globals.Warn(fmt.Sprintf("Error loading remote config: %v", err))
		return
	}
	if remoteConfig == nil { // 防止远程配置加载失败时覆盖本地配置
		globals.Warn("Remote config load failed, using existing configuration.")
		return
	}

	configMutex.Lock()
	defer configMutex.Unlock()
	localModels := customVisionConfig.CustomVisionModels                    // 保留当前的本地模型
	mergedModels := append(localModels, remoteConfig.CustomVisionModels...) // 合并本地和远程模型
	mergedModels = uniqueStringSlice(mergedModels)                          // 去重

	customVisionConfig.CustomVisionModels = mergedModels
	lastConfigUpdateTime = time.Now()
	globals.Info(fmt.Sprintf("Successfully load and merge remote config, current config:%v, update time: %v", customVisionConfig, lastConfigUpdateTime))
}

// loadConfig  首次加载配置，优先加载远程配置，本地加载失败时使用远程
func loadConfig() {
	localConfig, localErr := loadConfigFromFile("config/vision_models.json")
	remoteConfig, remoteErr := loadConfigFromURL(remoteConfigURL)

	if localErr != nil {
		globals.Warn(fmt.Sprintf("Error loading local config: %v", localErr))
		if remoteErr != nil {
			globals.Warn(fmt.Sprintf("Error loading remote config: %v, using default empty config.", remoteErr))
			customVisionConfig = Config{} // 本地远程都失败，使用默认空配置
		} else {
			globals.Info("Local config load failed, using remote config.")
			customVisionConfig = *remoteConfig // 本地失败，使用远程配置
		}
	} else {
		if remoteErr != nil {
			globals.Warn(fmt.Sprintf("Error loading remote config: %v, using local config.", remoteErr))
			customVisionConfig = *localConfig // 远程失败，使用本地配置
		} else {
			// 本地和远程都加载成功，合并去重
			mergedModels := append(localConfig.CustomVisionModels, remoteConfig.CustomVisionModels...)
			mergedModels = uniqueStringSlice(mergedModels)
			customVisionConfig.CustomVisionModels = mergedModels
			globals.Info("Successfully loaded and merged local and remote configs.")
		}
	}

	lastConfigUpdateTime = time.Now()
	globals.Info(fmt.Sprintf("Initial config loaded, current config:%v, update time: %v", customVisionConfig, lastConfigUpdateTime))
}

// loadConfigFromFile 从本地文件加载配置
func loadConfigFromFile(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// loadConfigFromURL 从远程 URL 加载配置
func loadConfigFromURL(url string) (*Config, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// IsCustomVisionModel 检查模型是否是自定义视觉模型
func IsCustomVisionModel(model string) bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	for _, customModel := range customVisionConfig.CustomVisionModels {
		if customModel == model {
			return true
		}
	}
	return false
}

// uniqueStringSlice 对字符串切片去重
func uniqueStringSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
